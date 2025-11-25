package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	core *zap.Logger
}

func newZapLogger(cfg *config) Logger {
	encCfg := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		CallerKey:    "caller",
		MessageKey:   "msg",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	var enc = zapcore.NewConsoleEncoder(encCfg)
	//enc = zapcore.NewJSONEncoder(encCfg)

	core := zapcore.NewCore(
		enc,
		zapcore.AddSync(cfg.output),
		getZapLevel(cfg.level),
	)

	return &zapLogger{core: zap.New(core)}
}

func (z *zapLogger) Log(level Level, msg string, fields ...Field) {
	lvl := getZapLevel(level)
	if !z.core.Level().Enabled(lvl) {
		return
	}

	zapFields := z.getZapField(fields)
	z.core.Log(lvl, msg, zapFields...)
}

func (z *zapLogger) Logd(_ int, level Level, msg string, fields ...Field) {
	z.Log(level, msg, fields...)
}

func (z *zapLogger) With(fields ...Field) Logger {
	zapFields := z.getZapField(fields)
	return &zapLogger{core: z.core.With(zapFields...)}

}

func (z *zapLogger) WithCallerSkip(skip int) Logger {
	return &zapLogger{core: z.core.WithOptions(zap.AddCallerSkip(skip))}
}

func (z *zapLogger) getZapField(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.K, field.V)
	}
	return zapFields
}

func getZapLevel(level Level) zapcore.Level {
	return zapcore.Level(level - 1)
}
