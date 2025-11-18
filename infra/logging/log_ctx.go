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

type ctxConvertor func(ctx context.Context) context.Context

var ctxConvertors []ctxConvertor

// RegisterCtxConvertor 注册的context convertor，比如gin.Context
func RegisterCtxConvertor(cc ctxConvertor) {
	ctxConvertors = append(ctxConvertors, cc)
}

func WithContext(ctx context.Context, fields ...Field) *LogContext {
	if ctx == nil {
		ctx = context.TODO()
	}
	c := withContext(ctx)
	if len(fields) > 0 {
		c.Fields = append(c.Fields, fields...)
	}
	return c
}

// 从 context 获取或创建 LogContext
func withContext(ctx context.Context) *LogContext {
	if logCtx, ok := ctx.(*LogContext); ok {
		return logCtx
	}

	for _, convertor := range ctxConvertors {
		if c := convertor(ctx); c != nil {
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
	withContext(ctx).log(DebugLevel, msg, fields...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	withContext(ctx).log(DebugLevel, fmt.Sprintf(format, args...))
}
func Info(ctx context.Context, msg string, fields ...Field) {
	withContext(ctx).log(InfoLevel, msg, fields...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	withContext(ctx).log(InfoLevel, fmt.Sprintf(format, args...))
}

func Error(ctx context.Context, msg string, fields ...Field) {
	withContext(ctx).log(ErrorLevel, msg, fields...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	withContext(ctx).log(ErrorLevel, fmt.Sprintf(format, args...))
}

func Errorw(ctx context.Context, err error, fields ...Field) {
	withContext(ctx).log(ErrorLevel, err.Error(), fields...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	withContext(ctx).log(WarnLevel, msg, fields...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	withContext(ctx).log(WarnLevel, fmt.Sprintf(format, args...))
}
