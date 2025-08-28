package api

import (
	"go-tpl/internal/api/handler"
	"go-tpl/internal/domain/repo"
	"go-tpl/internal/infra/global"
	"go-tpl/internal/service"

	"github.com/gin-gonic/gin"
)

type router struct {
	bidding *handler.BiddingHandler
}

func initAppContext() *router {
	// init repo
	var (
		ruleRepo        = repo.NewBiddingCustomRuleRepo(global.DB)
		otherRuleRepo   = repo.NewBiddingOtherRuleRepo(global.DB)
		orderRepo       = repo.NewOrderRepo(global.DB)
		orderExtendRepo = repo.NewOrderExtendRepo(global.DB)
		userRepo        = repo.NewUserRepo(global.DB)
		brandRepo       = repo.NewCinemaBrandRepo(global.DB)
		categoryRepo    = repo.NewCustomCinemaCategoryRepo(global.DB)
		biddingLogRepo  = repo.NewOrderBiddingLogRepo(global.DB)
		billRepo        = repo.NewBiddingBillRepo(global.DB)
		billBottomRepo  = repo.NewBiddingBillBottomRepo(global.DB)
		rulePriceRepo   = repo.NewBiddingRulePriceRepo(global.DB)
		oupAppRepo      = repo.NewOupAppRepo(global.DB)
	)

	// init service
	var (
		orderSvc   = service.NewOrderService(orderRepo, orderExtendRepo, global.Redis)
		biddingSvc = service.NewBiddingService(userRepo, ruleRepo, otherRuleRepo, brandRepo, categoryRepo,
			biddingLogRepo, billRepo, billBottomRepo, rulePriceRepo, oupAppRepo, orderSvc, global.Redis)
	)

	return &router{
		bidding: handler.NewBiddingHandler(biddingSvc),
	}
}

func SetupRoute(app *gin.Engine) {

	route := initAppContext()

	bid := app.Group("/bidding")
	{
		bid.GET("/list", route.bidding.ListWithEffective)
		bid.POST("/auto", route.bidding.AutoBidding)
	}
}
