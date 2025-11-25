package logger

import (
	"context"
	"errors"
	"testing"
)

func TestCtxInfo(t *testing.T) {
	Init()

	t.Run("CtxLog", func(t *testing.T) {
		ctx := context.Background()
		Info(ctx, "test")

		ctx = WithCtx(ctx, Str("requestId", "23123"))

		Warn(ctx, "test", Int("age", 10))
		Errw(ctx, errors.New("test error"), Int("age", 18))

		callerLog := Ctx(ctx)
		ctx = ToCtx(ctx, callerLog)
		Error(ctx, "test")

		callerLog.Log(INFO, "test", Str("requestId", "23123"))
		callerLog.Logd(0, WARN, "test", Str("requestId", "23123"))
	})
}
