package logx

import (
	"context"
	"go-tpl/infra/logger"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	GinHeaderLogID = "X-Trace-LogID"
)

var (
	ctxKeyLogID struct{}
	ctxKeyLog   struct{}

	once  sync.Once
	sugar *zap.SugaredLogger
)

func getLogger() *zap.SugaredLogger {
	once.Do(func() {
		sugar = logger.GetLogger().Sugar()
	})
	return sugar
}

// 从任意 Context 获取 Logger
func getLoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return getLogger()
	}

	// 如果是 Gin context
	if gCtx, gOk := ctx.(*gin.Context); gOk {
		if t, exist := gCtx.Get(ctxKeyLog); exist {
			if s, ok := t.(*zap.SugaredLogger); ok {
				return s
			}
		}
		log := getLogger().With(zap.String("logId", getLogIDFromGin(gCtx)))
		gCtx.Set(ctxKeyLog, log)
		return log
	}

	// 标准 context - 没有缓存，创建新的
	return getLogger().With(zap.String("logId", getLogIDFromContext(ctx)))
}

func getLogIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return generateLogID()
	}

	if logID, ok := ctx.Value(ctxKeyLogID).(string); ok {
		return logID
	}

	return generateLogID()
}

func getLogIDFromGin(ctx *gin.Context) string {
	if ctx == nil {
		return generateLogID()
	}
	if logId := ctx.GetString(ctxKeyLogID); logId != "" {
		return logId
	}
	// 尝试从header中获取
	var logId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		logId = ctx.GetHeader(GinHeaderLogID)
	}
	if logId == "" {
		logId = generateLogID()
	}
	ctx.Set(ctxKeyLogID, logId)
	return logId
}

func generateLogID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}

// 统一的日志函数

func WithLogID(ctx context.Context) context.Context {
	logID := getLogIDFromContext(ctx)
	return context.WithValue(ctx, ctxKeyLogID, logID)
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
