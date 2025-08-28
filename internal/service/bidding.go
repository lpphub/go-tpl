package service

import (
	"context"
	"fmt"
	"go-tpl/internal/api/vo"
	"go-tpl/internal/common/constant"
	"go-tpl/internal/common/dto"
	"go-tpl/internal/domain/entity"
	"go-tpl/internal/domain/repo"
	"go-tpl/internal/infra/global"
	"go-tpl/internal/service/checker"
	"go-tpl/pkg/ext"
	"go-tpl/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/logger/logx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BiddingService struct {
	userRepo       *repo.UserRepo
	ruleRepo       *repo.BiddingCustomRuleRepo
	otherRuleRepo  *repo.BiddingOtherRuleRepo
	brandRepo      *repo.CinemaBrandRepo
	categoryRepo   *repo.CustomCinemaCategoryRepo
	biddingLogRepo *repo.OrderBiddingLogRepo
	billRepo       *repo.BiddingBillRepo
	billBottomRepo *repo.BiddingBillBottomRepo
	rulePriceRepo  *repo.BiddingRulePriceRepo
	oupAppRepo     *repo.OupAppRepo
	orderSvc       *OrderService
	redis          *redis.Client

	// 批处理器
	processor *ext.AsyncProcessor[*dto.BiddingParamsDTO, []*dto.BiddingResultDTO]
}

func NewBiddingService(
	user *repo.UserRepo,
	rule *repo.BiddingCustomRuleRepo,
	otherRule *repo.BiddingOtherRuleRepo,
	brand *repo.CinemaBrandRepo,
	category *repo.CustomCinemaCategoryRepo,
	biddingLog *repo.OrderBiddingLogRepo,
	bill *repo.BiddingBillRepo,
	billBottom *repo.BiddingBillBottomRepo,
	rulePrice *repo.BiddingRulePriceRepo,
	oupApp *repo.OupAppRepo,
	orderSvc *OrderService,
	redis *redis.Client,
) *BiddingService {
	svc := &BiddingService{
		userRepo:       user,
		ruleRepo:       rule,
		otherRuleRepo:  otherRule,
		brandRepo:      brand,
		categoryRepo:   category,
		biddingLogRepo: biddingLog,
		billRepo:       bill,
		billBottomRepo: billBottom,
		rulePriceRepo:  rulePrice,
		oupAppRepo:     oupApp,
		orderSvc:       orderSvc,
		redis:          redis,
	}

	svc.processor, _ = ext.NewAsyncProcessor(
		svc.processFunc,
		ext.WithMaxConcurrency(2000),
	)
	return svc
}

func (s *BiddingService) ListWithEffective(ctx *gin.Context) ([]*entity.BiddingCustomRule, error) {
	return s.ruleRepo.ListWithEffective(ctx)
}

func (s *BiddingService) AutoBidding(ctx *gin.Context, params vo.BiddingAutoDTO) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logx.Infof(ctx, "自动报价耗时: %v\n", duration)
	}()

	order, err := s.orderSvc.GetOneForBidding(ctx, params)
	if err != nil {
		return errors.WithMessage(err, "未查询到订单")
	}
	logx.Infof(ctx, "订单信息: %d", order.Id)

	// 预加载数据
	_ = checker.NewLocalCacheLoader(s.brandRepo, s.categoryRepo).PreloadCache(ctx)

	customRuleList, _ := s.checkCustomRule(ctx, order, params)
	otherRuleList, _ := s.checkOtherRule(ctx, order, params)
	merged := append(customRuleList, otherRuleList...)
	logx.Infof(ctx, "匹配命中的规则总数: %d", len(merged))

	if len(merged) > 0 {
		biddingResult, err := s.rankingByUser(ctx, order, merged, params.IsLimitPrice)
		if err != nil {
			logx.Err(ctx, err, "ranking排序失败")
			return err
		}

		resultStr, _ := jsoniter.MarshalToString(biddingResult)
		logx.Infof(ctx, "自动竞价结果: %s", resultStr)

		if len(biddingResult) > 0 {
			//for _, bidding := range biddingResult {
			//if params.IsBidding {
			//	var (
			//		ruleId   int64
			//		ruleName string
			//		ruleStr  string
			//		ruleType = bidding.Type
			//	)
			//	if bidding.RuleInfo != nil {
			//		ruleId = bidding.RuleInfo.Id
			//		ruleName = bidding.RuleInfo.Name
			//		ruleStr, _ = jsoniter.MarshalToString(bidding.RuleInfo)
			//	} else if bidding.OtherRuleInfo != nil {
			//		ruleId = bidding.OtherRuleInfo.Id
			//		ruleName = bidding.OtherRuleInfo.Name
			//		ruleStr, _ = jsoniter.MarshalToString(bidding.OtherRuleInfo)
			//	}
			//	log := entity.OrderAutoBiddingLog{
			//		OrderId:    order.Id,
			//		Uid:        bidding.UID,
			//		RuleId:     ruleId,
			//		RuleName:   ruleName,
			//		RuleType:   ruleType,
			//		RuleInfo:   ruleStr,
			//		CreateTime: ext.NowTimestamp(),
			//	}
			//	_ = s.biddingLogRepo.Create(ctx, &log)
			//}
			//
			//// 写库
			//if err = s.processStoreData(ctx, order, bidding); err != nil {
			//	logx.Err(ctx, err, "保存数据失败")
			//	return err
			//}
			//
			//if params.IsSaveRedis {
			//	bidding.RuleInfo = nil
			//	biddingStr, _ := jsoniter.MarshalToString(bidding)
			//	s.redis.SAdd(ctx, fmt.Sprintf(constant.CacheBiddingResult, order.Id), biddingStr)
			//}
			//}
		}
	}
	return nil
}

func (s *BiddingService) processFunc(c context.Context, params *dto.BiddingParamsDTO) ([]*dto.BiddingResultDTO, error) {
	ctx := c.(*gin.Context)
	okList := make([]*dto.BiddingResultDTO, 0)

	if len(params.CustomRuleList) > 0 {
		ruleChecker := checker.NewCustomRuleChecker(params.Order)

		for _, rule := range params.CustomRuleList {
			checkResult, _ := ruleChecker.CheckRule(ctx, rule)
			if checkResult != nil && checkResult.Status == 1 {
				okList = append(okList, checkResult)
			}
		}
	}

	if len(params.OtherRuleList) > 0 {
		ruleChecker := checker.NewOtherRuleChecker(params.Order, s.rulePriceRepo)

		for _, rule := range params.OtherRuleList {
			checkResult, _ := ruleChecker.CheckRule(ctx, rule)
			if checkResult != nil && checkResult.Status == 1 {
				okList = append(okList, checkResult)
			}
		}
	}

	if len(okList) > 0 { // 按规则校验成功后，再按用户维度的条件校验
		logx.Infof(ctx, "批次规则匹配数量: %d", len(okList))

		userChecker := checker.NewUserChecker(params.Order, s.userRepo, s.billRepo, s.redis)
		return userChecker.Check(ctx, okList, params.IsBidding)
	}
	logx.Error(ctx, "批次规则未匹配到适用项")
	return nil, nil
}

func (s *BiddingService) checkCustomRule(ctx *gin.Context, order *dto.OrderDTO, params vo.BiddingAutoDTO) ([]*dto.BiddingResultDTO, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logx.Infof(ctx, "自定义规则耗时: %v\n", duration)
	}()

	ruleList, err := s.ruleRepo.ListWithEffective(ctx)
	if err != nil {
		return nil, err
	}
	logx.Infof(ctx, "有效自定义规则数量: %d", len(ruleList))

	ruleSlice, _ := util.Partition(ruleList, 200)
	var taskList []*dto.BiddingParamsDTO
	for _, slice := range ruleSlice {
		param := &dto.BiddingParamsDTO{
			Order:          order,
			CustomRuleList: slice,
			IsBidding:      params.IsBidding,
			IsLimitPrice:   params.IsLimitPrice,
		}
		taskList = append(taskList, param)
	}

	_ctx := ctx.Copy()
	resultList, err := s.processor.Process(_ctx, taskList)
	if err != nil {
		logx.Err(ctx, err, "自定义规则批处理有错误")
	}
	result := util.FlattenSlice(resultList)
	logx.Infof(ctx, "自定义规则匹配结果: %d", len(result))

	return result, nil
}

func (s *BiddingService) checkOtherRule(ctx *gin.Context, order *dto.OrderDTO, params vo.BiddingAutoDTO) ([]*dto.BiddingResultDTO, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logx.Infof(ctx, "其他规则耗时: %v\n", duration)
	}()

	ruleList, err := s.otherRuleRepo.ListWithEffective(ctx)
	if err != nil {
		return nil, err
	}
	logx.Infof(ctx, "有效其他规则数量: %d", len(ruleList))

	ruleIds := make([]int64, 0, len(ruleList))
	for _, rule := range ruleList {
		ruleIds = append(ruleIds, rule.Id)
	}
	priceList, err := s.rulePriceRepo.ListByRuleIds(ctx, ruleIds)
	if err != nil {
		return nil, err
	}
	rulePriceGroup := make(map[int64][]*entity.BiddingRulePrice)
	for _, price := range priceList {
		rulePriceGroup[price.RuleId] = append(rulePriceGroup[price.RuleId], price)
	}
	ruleDTOList := make([]*dto.RuleDTO, 0, len(ruleList))
	for _, rule := range ruleList {
		ruleDTOList = append(ruleDTOList, &dto.RuleDTO{
			OtherRule:      rule,
			OtherRulePrice: rulePriceGroup[rule.Id],
		})
	}

	ruleSlice, _ := util.Partition(ruleDTOList, 50)
	var taskList []*dto.BiddingParamsDTO
	for _, slice := range ruleSlice {
		param := &dto.BiddingParamsDTO{
			Order:         order,
			OtherRuleList: slice,
			IsBidding:     params.IsBidding,
			IsLimitPrice:  params.IsLimitPrice,
		}
		taskList = append(taskList, param)
	}

	resultList, err := s.processor.Process(ctx, taskList)
	if err != nil {
		logx.Err(ctx, err, "其他规则批处理有错误")
	}
	result := util.FlattenSlice(resultList)
	logx.Infof(ctx, "其他规则匹配结果: %d", len(result))
	return result, nil
}

func (s *BiddingService) rankingByUser(ctx *gin.Context, order *dto.OrderDTO, biddingResult []*dto.BiddingResultDTO, isLimitPrice bool) ([]*dto.BiddingResultDTO, error) {
	userRuleMap := make(map[int64][]*dto.BiddingResultDTO)
	for _, rule := range biddingResult {
		userRuleMap[rule.UID] = append(userRuleMap[rule.UID], rule)
	}

	result := make([]*dto.BiddingResultDTO, 0, len(biddingResult))
	for _, v := range userRuleMap {
		r, _ := s.processRanking(ctx, order, v, isLimitPrice)
		result = append(result, r)
	}
	return result, nil
}

func (s *BiddingService) processRanking(ctx *gin.Context, order *dto.OrderDTO, biddingResult []*dto.BiddingResultDTO, isLimitPrice bool) (*dto.BiddingResultDTO, error) {
	var bidPrice1, bidPrice2 *dto.BiddingResultDTO

	minPrice := decimal.Zero
	maxPrice := order.MaoyanPrice.Mul(order.Max).Round(1)
	percent95 := order.MaoyanPrice.Mul(decimal.NewFromFloat(0.95)).Round(1)
	newMaxPrice := maxPrice.Add(decimal.NewFromFloat(1.5))
	if newMaxPrice.GreaterThan(percent95) {
		newMaxPrice = percent95
	}

	for _, item := range biddingResult {
		if item.Status != 1 {
			continue
		}
		var priceStatus1, priceStatus2 bool
		if isLimitPrice {
			priceStatus1 = item.Price.GreaterThanOrEqual(minPrice) && item.Price.LessThanOrEqual(maxPrice)
			priceStatus2 = item.Price.GreaterThanOrEqual(minPrice) && item.Price.LessThanOrEqual(newMaxPrice)
		} else {
			priceStatus1 = true
			priceStatus2 = false
		}

		if priceStatus1 {
			if bidPrice1 == nil || item.Price.LessThan(bidPrice1.Price) {
				bidPrice1 = item
			}
		} else if priceStatus2 {
			if bidPrice2 == nil || item.Price.LessThan(bidPrice2.Price) {
				bidPrice2 = item
			}
		}
	}

	if bidPrice1 == nil && bidPrice2 != nil {
		jsonStr, _ := jsoniter.MarshalToString(bidPrice2)
		logx.Infof(ctx, "使用新价格: %s", jsonStr)
		//s.redis.SAdd(ctx, fmt.Sprintf(constant.CacheBiddingNewMaxPrice, order.Id), jsonStr)
	}
	return bidPrice1, nil
}

func (s *BiddingService) processStoreData(ctx *gin.Context, order *dto.OrderDTO, bidding *dto.BiddingResultDTO) error {
	var (
		uid  = bidding.UID
		user = bidding.UserInfo
	)

	bill := entity.BiddingBill{
		AutoId:     bidding.RuleId,
		Uid:        uid,
		Price:      bidding.Price,
		OrderId:    order.Id,
		CinemaId:   order.CinemaId,
		IsSeat:     order.IsSeat,
		LoverSeat:  order.LoverSeat,
		IsAuto:     1,
		Status:     0,
		BindTime:   0,
		CreateTime: int(time.Now().Unix()),
	}

	// 检查是否为Oup用户
	isOupUser := s.oupAppRepo.Exist(ctx, uid)
	if isOupUser {
		bill.IsConfirm = 1
		bill.ConfirmTime = int(time.Now().Unix())
	}

	if order.Status != 0 {
		return errors.New("订单状态非竞价中")
	}

	// 开启事务
	err := global.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.billRepo.InsertWithTx(tx, &bill); err != nil {
			return err
		}

		key := fmt.Sprintf(constant.CacheBiddingUserForbidden, uid)
		if exists := s.redis.Exists(ctx, key).Val(); exists == 0 {
			updates := map[string]any{
				"forbidden_end_time": time.Now().AddDate(0, 0, 10).Format("2006-01-02 15:04:05"),
			}
			if err := s.userRepo.UpdateWithTx(tx, uid, updates); err != nil {
				return err
			}
			s.redis.SetEx(ctx, key, "1", 600*time.Second)
		}

		if user.IsAutoForbidden == 1 {
			updates := map[string]any{
				"forbidden": 0,
			}
			if err := s.userRepo.UpdateWithTx(tx, uid, updates); err != nil {
				return err
			}
		}

		if err := s.redis.SAdd(ctx, fmt.Sprintf(constant.CacheBiddingPay, order.Id, uid), bidding.Price).Err(); err != nil {
			return errors.WithMessage(err, "redis error")
		}

		if err := s.redis.SAdd(ctx, fmt.Sprintf(constant.CacheBiddingJoin, order.Id), uid).Err(); err != nil {
			return errors.WithMessage(err, "redis error")
		}

		if bidding.IfBottomCover {
			bottom := entity.BiddingBillBottomCover{
				BiddingId:  bill.Id,
				OrderId:    order.Id,
				Uid:        uid,
				CreateTime: ext.NowTimestamp(),
			}
			return s.billBottomRepo.InsertWithTx(tx, &bottom)
		}
		return nil
	})
	if err != nil {
		s.redis.SRem(ctx, fmt.Sprintf(constant.CacheBiddingPay, order.Id, uid), bidding.Price)
		s.redis.SRem(ctx, fmt.Sprintf(constant.CacheBiddingJoin, order.Id), uid)
		return err
	}
	return nil
}
