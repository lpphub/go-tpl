package handler

import (
	"go-tpl/internal/common/errs"
	"go-tpl/internal/service"
	"go-tpl/pkg/pagination"

	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
	"github.com/lpphub/golib/web"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) PageList(ctx *gin.Context) {
	var req pagination.Pagination
	if err := ctx.ShouldBind(&req); err != nil {
		logx.Error(ctx, err.Error())
		web.JsonWithError(ctx, errs.ErrParamInvalid)
		return
	}

	data, err := h.svc.PageList(ctx, req)
	if err != nil {
		web.JsonWithError(ctx, err)
	} else {
		web.JsonWithSuccess(ctx, data)
	}
}
