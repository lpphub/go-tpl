package logx

import (
	"go-tpl/infra/logging"

	"github.com/gin-gonic/gin"
)

const (
	GinHeaderLogID = "X-Trace-LogID"
)

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
		logId = ctx.GetHeader(GinHeaderLogID)
	}

	if logId == "" {
		logId = logging.GenerateLogID()
	}
	return logId
}
