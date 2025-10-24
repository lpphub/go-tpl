package logx

import (
	"context"
	"go-tpl/infra/logging"

	"github.com/gin-gonic/gin"
)

const (
	LogIDHeader = "X-Trace-LogID"
)

func init() {
	logging.RegisterCtxConvertor(func(ctx context.Context) context.Context {
		if gCtx, ok := ctx.(*gin.Context); ok {
			return gCtx.Request.Context()
		}
		return nil
	})
}
func GinLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logCtx := logging.WithLogID(c.Request.Context(), getLogIDFromGin(c))

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
		logId = logging.GenerateLogID()
	}
	return logId
}
