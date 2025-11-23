package logx

import (
	"fmt"
	"go-tpl/infra/logger"

	"github.com/gin-gonic/gin"
)

const (
	LogIDHeader = "X-Trace-LogID"
)

func GinLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logCtx := logger.WithField(c.Request.Context(), "log_id", getLogIDFromGin(c))

		logger.Info(logCtx, "gin request",
			logger.Str("path", fmt.Sprintf("[%s %s]", c.Request.Method, c.Request.RequestURI)))

		c.Request = c.Request.WithContext(logCtx)

		c.Next()
	}
}

func getLogIDFromGin(ctx *gin.Context) string {
	// 尝试从header中获取
	var logId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		logId = ctx.GetHeader(LogIDHeader)
	}

	if logId == "" {
		logId = logger.GenerateLogID()
	}
	return logId
}
