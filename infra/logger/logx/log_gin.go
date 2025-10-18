package logx

import (
	"context"
	"go-tpl/infra"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	GinHeaderLogId = "X-Trace-logID"
)

type (
	keyLogger struct{}
	keyLogID  struct{}
)

var (
	logGinKey keyLogger
	logIDKey  keyLogID

	once  sync.Once
	sugar *zap.SugaredLogger
)

func getLogger() *zap.SugaredLogger {
	once.Do(func() {
		sugar = infra.Logger.Sugar()
	})
	return sugar
}

func getLoggerFromContext(c context.Context) *zap.SugaredLogger {
	if gCtx, isOk := c.(*gin.Context); isOk {
		if t, exist := gCtx.Get(logGinKey); exist {
			if s, ok := t.(*zap.SugaredLogger); ok {
				return s
			}
		}
		sLog := sugar.With(zap.String("logId", getLogID(gCtx)))
		gCtx.Set(logGinKey, sLog)
		return sLog
	}

	return getLogger()
}

func getLogID(ctx *gin.Context) string {
	if ctx == nil {
		return generateLogID()
	}
	if logId := ctx.GetString(logIDKey); logId != "" {
		return logId
	}
	// 尝试从header中获取
	var logId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		logId = ctx.GetHeader(GinHeaderLogId)
	}
	if logId == "" {
		logId = generateLogID()
	}
	ctx.Set(logIDKey, logId)
	return logId
}

func generateLogID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}

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
