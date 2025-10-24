package logging

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LogContext struct {
	context.Context

	Logger Logger
	Fields []Field
}

func (l *LogContext) log(level LogLevel, msg string) {
	l.Logger.Write(level, msg, l.Fields...)
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

func GenerateLogID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}

func getLoggerFromContext(ctx context.Context) *LogContext {
	if ctx == nil {
		return WithContext(context.TODO())
	}
	// logCtx
	if logCtx, ok := ctx.(*LogContext); ok {
		return logCtx
	}

	// Gin context
	if gCtx, ok := ctx.(*gin.Context); ok {
		return getLoggerFromContext(gCtx.Request.Context())
	}

	return WithContext(ctx)
}

// 统一的日志函数

func Debug(ctx context.Context, msg string) {
	getLoggerFromContext(ctx).log(DebugLevel, msg)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	Debug(ctx, fmt.Sprintf(format, args))

}
func Info(ctx context.Context, msg string) {
	getLoggerFromContext(ctx).log(InfoLevel, msg)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	Info(ctx, fmt.Sprintf(format, args))
}

func Error(ctx context.Context, msg string) {
	getLoggerFromContext(ctx).log(ErrorLevel, msg)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	Error(ctx, fmt.Sprintf(format, args))
}

func Warn(ctx context.Context, msg string) {
	getLoggerFromContext(ctx).log(WarnLevel, msg)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	Warn(ctx, fmt.Sprintf(format, args))
}
