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

const (
	GinHeaderLogID = "X-Trace-LogId"
)

var (
	ctxKeyLogID  struct{}
	ctxKeyLogger struct{}

	once  sync.Once
	sugar *zap.SugaredLogger
)

type LogContext interface {
	context.Context

	LogID() string
	Logger() *zap.SugaredLogger
	WithValue(key, value interface{}) LogContext
}

type logContext struct {
	context.Context
	logID  string
	logger *zap.SugaredLogger
}

func (c *logContext) LogID() string {
	return c.logID
}

func (c *logContext) Logger() *zap.SugaredLogger {
	return c.logger
}

func (c *logContext) WithValue(key, value interface{}) LogContext {
	return &logContext{
		Context: context.WithValue(c.Context, key, value),
		logID:   c.logID,
		logger:  c.logger,
	}
}

// NewLogContext 从标准 context 创建 LogContext
func NewLogContext(ctx context.Context) LogContext {
	if ctx == nil {
		ctx = context.Background()
	}

	logID := extractLogID(ctx)

	log := getLogger().With(zap.String("logId", logID))

	return &logContext{
		Context: ctx,
		logID:   logID,
		logger:  log,
	}
}

// NewLogContextFromGin 从 Gin context 创建 LogContext
func NewLogContextFromGin(ginCtx *gin.Context) LogContext {
	if ginCtx == nil {
		return NewLogContext(nil)
	}

	logID := extractLogIDFromGin(ginCtx)

	log := getLogger().With(zap.String("logId", logID))

	// 使用 Gin context 作为底层 context
	return &logContext{
		Context: ginCtx,
		logID:   logID,
		logger:  log,
	}
}

// WithLogID 创建带有指定 LogID 的新 LogContext
func WithLogID(ctx context.Context, logID string) LogContext {
	if ctx == nil {
		ctx = context.Background()
	}

	log := getLogger().With(zap.String("logId", logID))

	return &logContext{
		Context: ctx,
		logID:   logID,
		logger:  log,
	}
}

// getLogContext 获取或创建 LogContext
func getLogContext(ctx context.Context) LogContext {
	if ctx == nil {
		return NewLogContext(nil)
	}

	// 如果已经是 LogContext，直接返回
	if logCtx, ok := ctx.(LogContext); ok {
		return logCtx
	}

	// 如果是 Gin context，使用专门的创建函数
	if _, ok := ctx.(*gin.Context); ok {
		return NewLogContextFromGin(ctx.(*gin.Context))
	}

	// 标准 context
	return NewLogContext(ctx)
}

// extractLogID 从标准 context 中提取 LogID
func extractLogID(ctx context.Context) string {
	if ctx == nil {
		return generateLogID()
	}

	// 从 context 值中获取
	if logID, ok := ctx.Value(ctxKeyLogID).(string); ok {
		return logID
	}

	return generateLogID()
}

// extractLogIDFromGin 从 Gin context 中提取 LogID
func extractLogIDFromGin(ginCtx *gin.Context) string {
	if ginCtx == nil {
		return generateLogID()
	}

	// 从 Header 中获取
	if ginCtx.Request != nil && ginCtx.Request.Header != nil {
		if logID := ginCtx.GetHeader(GinHeaderLogID); logID != "" {
			return logID
		}
	}

	// 从 Gin context 存储中获取
	logID := ginCtx.GetString(ctxKeyLogID)
	if logID != "" {
		return logID
	}
	return generateLogID()
}

func getLogger() *zap.SugaredLogger {
	once.Do(func() {
		sugar = logger.GetLogger().Sugar()
	})
	return sugar
}

// generateLogID 生成 LogID
func generateLogID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}

// 统一的日志函数

func Info(ctx context.Context, msg string) {
	getLogContext(ctx).Logger().Info(msg)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	getLogContext(ctx).Logger().Infof(format, args...)
}

func Error(ctx context.Context, msg string) {
	getLogContext(ctx).Logger().Error(msg)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	getLogContext(ctx).Logger().Errorf(format, args...)
}

func Debug(ctx context.Context, msg string) {
	getLogContext(ctx).Logger().Debug(msg)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	getLogContext(ctx).Logger().Debugf(format, args...)
}

func Warn(ctx context.Context, msg string) {
	getLogContext(ctx).Logger().Warn(msg)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	getLogContext(ctx).Logger().Warnf(format, args...)
}
