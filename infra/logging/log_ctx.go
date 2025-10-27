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

type ctxConvertor func(ctx context.Context) context.Context

var ctxConvertors []ctxConvertor

// RegisterCtxConvertor 注册的context convertor，比如gin.Context
func RegisterCtxConvertor(cc ctxConvertor) {
	ctxConvertors = append(ctxConvertors, cc)
}

func WithContext(ctx context.Context) *LogContext {
	if logCtx, ok := ctx.(*LogContext); ok {
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

	for _, fn := range ctxConvertors {
		if c := fn(ctx); c != nil {
			ctx = c
		}
	}

	return WithContext(ctx)
}

// 统一的日志函数

func Debug(ctx context.Context, msg string) {
	getLoggerFromCtx(ctx).log(DebugLevel, msg)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	getLoggerFromCtx(ctx).log(DebugLevel, fmt.Sprintf(format, args))
}
func Info(ctx context.Context, msg string) {
	getLoggerFromCtx(ctx).log(InfoLevel, msg)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	getLoggerFromCtx(ctx).log(InfoLevel, fmt.Sprintf(format, args))
}

func Error(ctx context.Context, msg string) {
	getLoggerFromCtx(ctx).log(ErrorLevel, msg)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	getLoggerFromCtx(ctx).log(ErrorLevel, fmt.Sprintf(format, args))
}

func Errorw(ctx context.Context, err error) {
	getLoggerFromCtx(ctx).log(ErrorLevel, err.Error())
}

func Warn(ctx context.Context, msg string) {
	getLoggerFromCtx(ctx).log(WarnLevel, msg)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	getLoggerFromCtx(ctx).log(WarnLevel, fmt.Sprintf(format, args))
}
