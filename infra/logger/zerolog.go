package logger

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type zeroLogger struct{ core zerolog.Logger }

func newZeroLogger(cfg *config) Logger {
	// default config
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	output := os.Stdout

	l := zerolog.New(output).With().Timestamp().Logger().Level(zerolog.Level(cfg.level))
	return &zeroLogger{core: l}
}

func (z *zeroLogger) Log(level Level, msg string, fields ...Field) {
	z.Logd(1, level, msg, fields...)
}

func (z *zeroLogger) Logd(depth int, level Level, msg string, fields ...Field) {
	e := z.event(level)

	// skip frame
	e = e.CallerSkipFrame(depth + 2)

	for _, f := range fields {
		e = z.addField(e, f)
	}
	e.Msg(msg)
}

func (z *zeroLogger) With(fields ...Field) Logger {
	ctx := z.core.With()
	for _, f := range fields {
		ctx = z.withField(ctx, f)
	}
	return &zeroLogger{core: ctx.Logger()}
}

func (z *zeroLogger) WithCallerSkip(skip int) Logger {
	return &zeroLogger{core: z.core.With().CallerWithSkipFrameCount(skip).Logger()}
}

func (z *zeroLogger) event(level Level) *zerolog.Event {
	switch level {
	case DEBUG:
		return z.core.Debug()
	case WARN:
		return z.core.Warn()
	case ERROR:
		return z.core.Error()
	case FATAL:
		return z.core.Fatal()
	default:
		return z.core.Info()
	}
}

func (z *zeroLogger) addField(e *zerolog.Event, f Field) *zerolog.Event {
	switch v := f.V.(type) {
	case string:
		return e.Str(f.K, v)
	case int:
		return e.Int(f.K, v)
	case int64:
		return e.Int64(f.K, v)
	case bool:
		return e.Bool(f.K, v)
	case error:
		return e.AnErr(f.K, v)
	default:
		return e.Interface(f.K, v)
	}
}

func (z *zeroLogger) withField(c zerolog.Context, f Field) zerolog.Context {
	switch v := f.V.(type) {
	case string:
		return c.Str(f.K, v)
	case int:
		return c.Int(f.K, v)
	case int64:
		return c.Int64(f.K, v)
	case bool:
		return c.Bool(f.K, v)
	case error:
		return c.AnErr(f.K, v)
	default:
		return c.Interface(f.K, v)
	}
}
