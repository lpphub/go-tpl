package zlog

import (
	"context"
	"testing"
)

func TestInit(t *testing.T) {
	Init(nil)

	ctx := WithField(context.Background(), "logId", "123213")

	Info(ctx, "hello world")

	Error(ctx, "hello world2", Str("test", "ok"))

}
