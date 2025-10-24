package logging

import (
	"context"
	"fmt"
)

type LogContext struct {
	context.Context

	Logger Logger
	Fields []Field
}

func (l *LogContext) log(level LogLevel, msg string) {
	l.Logger.Write(level, msg, l.Fields...)
}

type logCtxExtractor func(ctx context.Context) context.Context

var ctxExtractors []logCtxExtractor

func RegisterCtxExtractor(ce func(ctx context.Context) context.Context) {
	ctxExtractors = append(ctxExtractors, ce)
}

func WithContext(ctx context.Context) *LogContext {
	logCtx, ok := ctx.(*LogContext)
	if ok {
		return logCtx
	}
	return &LogContext{
		Context: ctx,
		Logger:  GetLogger(),
	}
}

func WithLogID(ctx context.Context, logID string) *LogContext {
	logCtx := WithContext(ctx)

	logCtx.Fields = append(logCtx.Fields, Field{Key: "logID", Value: logID})
	return logCtx
}

func getLoggerFromCtx(ctx context.Context) *LogContext {
	if ctx == nil {
		return WithContext(context.TODO())
	}

	// 兼容注册的context extractor，比如gin.Context
	for _, fn := range ctxExtractors {
		if gCtx := fn(ctx); gCtx != nil {
			ctx = gCtx
		}
	}

	if logCtx, ok := ctx.(*LogContext); ok {
		return logCtx
	}

	return WithContext(ctx)
}

// 统一的日志函数

func Debug(ctx context.Context, msg string) {
	getLoggerFromCtx(ctx).log(DebugLevel, msg)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	Debug(ctx, fmt.Sprintf(format, args))

}
func Info(ctx context.Context, msg string) {
	getLoggerFromCtx(ctx).log(InfoLevel, msg)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	Info(ctx, fmt.Sprintf(format, args))
}

func Error(ctx context.Context, msg string) {
	getLoggerFromCtx(ctx).log(ErrorLevel, msg)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	Error(ctx, fmt.Sprintf(format, args))
}

func Warn(ctx context.Context, msg string) {
	getLoggerFromCtx(ctx).log(WarnLevel, msg)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	Warn(ctx, fmt.Sprintf(format, args))
}
