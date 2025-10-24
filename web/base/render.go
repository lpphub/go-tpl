package base

import (
	"errors"
	"go-tpl/logic/shared"
	"go-tpl/web/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, types.Resp{
		Code: 0,
		Msg:  "ok",
	})
}

func OKWithData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, types.Resp{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

func Fail(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(http.StatusOK, types.Resp{
		Code: code,
		Msg:  msg,
	})
}

func FailWithErr(ctx *gin.Context, code int, err error) {
	ctx.AbortWithStatusJSON(http.StatusOK, types.Resp{
		Code: code,
		Msg:  err.Error(),
	})
}

func FailWithError(ctx *gin.Context, err error) {
	var e shared.Error
	if ok := errors.As(err, &e); ok {
		Fail(ctx, e.Code, e.Error())
	} else {
		FailWithErr(ctx, shared.ErrServerError.Code, err)
	}
}

func FailWithStatus(ctx *gin.Context, statusCode int, err error) {
	ctx.AbortWithStatusJSON(statusCode, types.Resp{
		Code: statusCode,
		Msg:  err.Error(),
	})
}
