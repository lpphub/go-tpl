package logx

import (
	"context"
	"fmt"
	"go-tpl/infra/logging"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	LogIDHeader = "X-Trace-LogID"
)

func init() {
	logging.RegisterCtxAdapter(func(ctx context.Context) context.Context {
		if gCtx, ok := ctx.(*gin.Context); ok {
			return gCtx.Request.Context()
		}
		return ctx
	})
}

func GinLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logCtx := logging.WithContext(c.Request.Context(), logging.WithField("log_id", getLogIDFromGin(c)))

		logging.Info(logCtx, "gin request",
			logging.WithField("path", fmt.Sprintf("[%s %s]", c.Request.Method, c.Request.RequestURI)))

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
		logId = GenerateLogID()
	}
	return logId
}

func GenerateLogID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}
