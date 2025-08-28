package ext

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type Processor[T any, R any] interface {
	Process(ctx context.Context, items []T) ([]R, error)
	ProcessWithTimeout(ctx context.Context, items []T, timeout time.Duration) ([]R, error)
}

type config struct {
	MaxConcurrency int
	DefaultTimeout time.Duration
}

type Option func(*config)

type AsyncProcessor[T any, R any] struct {
	cfg       config
	processFn func(context.Context, T) (R, error)
}

type result[R any] struct {
	index int
	res   R
	err   error
}

func defaultConfig() config {
	return config{
		MaxConcurrency: 10,
		DefaultTimeout: 30 * time.Second,
	}
}

func NewAsyncProcessor[T any, R any](
	processFn func(context.Context, T) (R, error),
	opts ...Option,
) (*AsyncProcessor[T, R], error) {
	if processFn == nil {
		return nil, errors.New("process function cannot be nil")
	}

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.MaxConcurrency <= 0 {
		return nil, errors.New("max concurrency must be greater than 0")
	}

	return &AsyncProcessor[T, R]{
		cfg:       cfg,
		processFn: processFn,
	}, nil
}

func WithMaxConcurrency(max int) Option {
	return func(c *config) {
		if max > 0 {
			c.MaxConcurrency = max
		}
	}
}

func WithDefaultTimeout(timeout time.Duration) Option {
	return func(c *config) {
		if timeout > 0 {
			c.DefaultTimeout = timeout
		}
	}
}

func (p *AsyncProcessor[T, R]) Process(ctx context.Context, items []T) ([]R, error) {
	return p.processWithCtx(ctx, items)
}

func (p *AsyncProcessor[T, R]) ProcessWithTimeout(ctx context.Context, items []T, timeout time.Duration) ([]R, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return p.processWithCtx(ctx, items)
}

func (p *AsyncProcessor[T, R]) processWithCtx(ctx context.Context, items []T) ([]R, error) {
	if len(items) == 0 {
		return []R{}, nil
	}

	resultCh := make(chan result[R], len(items))

	g, _c := errgroup.WithContext(ctx)
	g.SetLimit(p.cfg.MaxConcurrency)

	for i, item := range items {
		idx, t := i, item
		g.Go(func() error {
			select {
			case <-_c.Done():
				return _c.Err()
			default:
				res, err := p.processFn(ctx, t)
				resultCh <- result[R]{
					index: idx,
					res:   res,
					err:   err,
				}
				return nil
			}
		})
	}

	go func() {
		_ = g.Wait()
		close(resultCh)
	}()

	results := make([]R, 0, len(items))
	var errs []error
	for r := range resultCh {
		if r.err != nil {
			errs = append(errs, fmt.Errorf("%d: %w", r.index+1, r.err))
		} else {
			results = append(results, r.res)
		}
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return results, &BatchError{Errors: errs}
	}
	return results, nil
}

type BatchError struct {
	Errors []error
}

func (e *BatchError) Error() string {
	errMsg := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		errMsg[i] = err.Error()
	}
	return fmt.Sprintf("processed errors: %s", strings.Join(errMsg, ","))
	//return fmt.Sprintf("processed errors count: %d", len(e.Errors))
}

func (e *BatchError) AllErrors() []error {
	return append([]error(nil), e.Errors...)
}
