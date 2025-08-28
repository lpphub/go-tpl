package handler

import (
	"go-tpl/internal/api/vo"
	"go-tpl/internal/infra/errs"
	"go-tpl/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
	"github.com/lpphub/golib/web"
)

type BiddingHandler struct {
	svc *service.BiddingService
}

func NewBiddingHandler(svc *service.BiddingService) *BiddingHandler {
	return &BiddingHandler{
		svc: svc,
	}
}

func (h *BiddingHandler) ListWithEffective(ctx *gin.Context) {
	data, err := h.svc.ListWithEffective(ctx)
	if err != nil {
		web.JsonWithError(ctx, err)
	} else {
		web.JsonWithSuccess(ctx, data)
	}
}

func (h *BiddingHandler) AutoBidding(ctx *gin.Context) {
	var req vo.BiddingAutoDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logx.Error(ctx, err.Error())
		web.JsonWithError(ctx, errs.ErrParamInvalid)
		return
	}

	err := h.svc.AutoBidding(ctx, req)
	if err != nil {
		web.JsonWithError(ctx, err)
	} else {
		web.JsonWithSuccess(ctx, "ok")
	}
}
