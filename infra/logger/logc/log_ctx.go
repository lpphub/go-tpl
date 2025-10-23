package logc

import (
	"context"
	"go-tpl/infra/logger"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type logContext struct {
	context.Context
	logID  string
	logger *zap.SugaredLogger
}

var (
	once  sync.Once
	sugar *zap.SugaredLogger
)

const (
	keyLogID = "logID"
)

func getLogger() *zap.SugaredLogger {
	once.Do(func() {
		sugar = logger.GetLogger().Sugar()
	})
	return sugar
}

// 从 Context 获取 Logger
func getLoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return getLogger()
	}
	// logCtx
	if logCtx, ok := ctx.(*logContext); ok {
		return logCtx.logger
	}

	// Gin context
	if gCtx, ok := ctx.(*gin.Context); ok {
		return getLoggerFromContext(gCtx.Request.Context())
	}

	return getLogger()
}

func generateLogID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}

// WithContext 创建带有logger的Context
func WithContext(ctx context.Context) context.Context {
	if ctx == nil {
		return WithContext(context.Background())
	}
	logCtx, ok := ctx.(*logContext)
	if ok {
		return logCtx
	}

	logID := generateLogID()

	return &logContext{
		Context: ctx,
		logID:   logID,
		logger:  getLogger().With(zap.String(keyLogID, logID)),
	}
}

func WithLogID(ctx context.Context, logID string) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	logCtx, ok := ctx.(*logContext)
	if ok {
		logCtx.logID = logID
		return logCtx
	}

	return &logContext{
		Context: ctx,
		logID:   logID,
		logger:  getLogger().With(zap.String(keyLogID, logID)),
	}
}

// 统一的日志函数

func Info(ctx context.Context, msg string) {
	getLoggerFromContext(ctx).Info(msg)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	getLoggerFromContext(ctx).Infof(format, args...)
}

func Error(ctx context.Context, msg string) {
	getLoggerFromContext(ctx).Error(msg)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	getLoggerFromContext(ctx).Errorf(format, args...)
}

func Debug(ctx context.Context, msg string) {
	getLoggerFromContext(ctx).Debug(msg)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	getLoggerFromContext(ctx).Debugf(format, args...)
}

func Warn(ctx context.Context, msg string) {
	getLoggerFromContext(ctx).Warn(msg)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	getLoggerFromContext(ctx).Warnf(format, args...)
}
