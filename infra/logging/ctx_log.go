package logging

import (
	"context"
	"fmt"
)

type LogContext struct {
	context.Context

	logger Logger
	Fields []Field
}

type ctxAdapter func(ctx context.Context) context.Context

var ctxAdapters []ctxAdapter

// RegisterCtxAdapter 注册的context adapter，比如gin.Context
func RegisterCtxAdapter(ca ctxAdapter) {
	ctxAdapters = append(ctxAdapters, ca)
}

func WithContext(ctx context.Context, fields ...Field) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}

	return &LogContext{
		Context: ctx,
		logger:  GetLogger(),
		Fields:  fields,
	}
}

func FromContext(ctx context.Context) *LogContext {
	if logCtx, ok := ctx.(*LogContext); ok {
		return logCtx
	}

	for _, adapter := range ctxAdapters {
		if c := adapter(ctx); c != nil {
			if logCtx, ok := c.(*LogContext); ok {
				return logCtx
			}
			return &LogContext{
				Context: c,
				logger:  GetLogger(),
			}
		}
	}

	return &LogContext{
		Context: ctx,
		logger:  GetLogger(),
	}
}

func (l *LogContext) log(level LogLevel, msg string, fields ...Field) {
	if len(fields) > 0 {
		fields = append(fields, l.Fields...)
	} else {
		fields = l.Fields
	}
	l.logger.Write(level, msg, fields...)
}

// 统一的日志函数

func Debug(ctx context.Context, msg string, fields ...Field) {
	FromContext(ctx).log(DebugLevel, msg, fields...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).log(DebugLevel, fmt.Sprintf(format, args...))
}
func Info(ctx context.Context, msg string, fields ...Field) {
	FromContext(ctx).log(InfoLevel, msg, fields...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).log(InfoLevel, fmt.Sprintf(format, args...))
}

func Error(ctx context.Context, msg string, fields ...Field) {
	FromContext(ctx).log(ErrorLevel, msg, fields...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).log(ErrorLevel, fmt.Sprintf(format, args...))
}

func Errorw(ctx context.Context, err error, fields ...Field) {
	FromContext(ctx).log(ErrorLevel, err.Error(), fields...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	FromContext(ctx).log(WarnLevel, msg, fields...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).log(WarnLevel, fmt.Sprintf(format, args...))
}
