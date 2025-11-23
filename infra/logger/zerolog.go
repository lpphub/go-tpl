// logger/zerolog.go
package logger

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	logger zerolog.Logger
}

func NewZeroLogger(cfg *Config) (*ZeroLogger, error) {
	w := cfg.Writer
	if w == nil {
		w = os.Stdout
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	logger := zerolog.New(w).
		With().
		Timestamp().
		CallerWithSkipFrameCount(4).
		Logger().
		Level(getLevel(cfg.Level))

	return &ZeroLogger{logger: logger}, nil
}

func getLevel(level Level) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *ZeroLogger) getLogger(ctx context.Context) *zerolog.Logger {
	if ctx == nil {
		return &l.logger
	}

	logger := zerolog.Ctx(ctx)
	if logger.GetLevel() == zerolog.Disabled {
		return &l.logger
	}

	return logger
}

func (l *ZeroLogger) applyFields(event *zerolog.Event, fields []Field) *zerolog.Event {
	for _, field := range fields {
		switch field.Type {
		case StringType:
			if v, ok := field.Value.(string); ok {
				event = event.Str(field.Key, v)
			}
		case IntType:
			if v, ok := field.Value.(int); ok {
				event = event.Int(field.Key, v)
			}
		case Int64Type:
			if v, ok := field.Value.(int64); ok {
				event = event.Int64(field.Key, v)
			}
		case Uint64Type:
			if v, ok := field.Value.(uint64); ok {
				event = event.Uint64(field.Key, v)
			}
		case BoolType:
			if v, ok := field.Value.(bool); ok {
				event = event.Bool(field.Key, v)
			}
		case Float64Type:
			if v, ok := field.Value.(float64); ok {
				event = event.Float64(field.Key, v)
			}
		case DurationType:
			if v, ok := field.Value.(time.Duration); ok {
				event = event.Dur(field.Key, v)
			}
		case TimeType:
			if v, ok := field.Value.(time.Time); ok {
				event = event.Time(field.Key, v)
			}
		case ErrorType:
			if v, ok := field.Value.(error); ok {
				event = event.AnErr(field.Key, v)
			}
		default:
			event = event.Interface(field.Key, field.Value)
		}
	}
	return event
}

func (l *ZeroLogger) getEvent(ctx context.Context, level Level) *zerolog.Event {
	logger := l.getLogger(ctx)
	switch level {
	case DebugLevel:
		return logger.Debug()
	case InfoLevel:
		return logger.Info()
	case WarnLevel:
		return logger.Warn()
	case ErrorLevel:
		return logger.Error()
	case FatalLevel:
		return logger.Fatal()
	default:
		return logger.Info()
	}
}

// Log 实现核心日志方法
func (l *ZeroLogger) Log(ctx context.Context, level Level, msg string, fields ...Field) {
	event := l.getEvent(ctx, level)
	l.applyFields(event, fields).Msg(msg)
}

// WithContext 将字段添加到 context
func (l *ZeroLogger) WithContext(ctx context.Context, fields ...Field) context.Context {
	logger := l.getLogger(ctx)
	logCtx := logger.With()

	for _, field := range fields {
		switch field.Type {
		case StringType:
			if v, ok := field.Value.(string); ok {
				logCtx = logCtx.Str(field.Key, v)
			}
		case IntType:
			if v, ok := field.Value.(int); ok {
				logCtx = logCtx.Int(field.Key, v)
			}
		case Int64Type:
			if v, ok := field.Value.(int64); ok {
				logCtx = logCtx.Int64(field.Key, v)
			}
		case BoolType:
			if v, ok := field.Value.(bool); ok {
				logCtx = logCtx.Bool(field.Key, v)
			}
		default:
			logCtx = logCtx.Interface(field.Key, field.Value)
		}
	}

	newLogger := logCtx.Logger()
	return newLogger.WithContext(ctx)
}

func (l *ZeroLogger) WithCaller(skip int) Logger {
	log := l.logger.With().CallerWithSkipFrameCount(skip).Logger()
	return &ZeroLogger{logger: log}
}
