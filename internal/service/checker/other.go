package checker

import (
	"errors"
	"go-tpl/internal/common/dto"
	"go-tpl/internal/domain/entity"
	"go-tpl/internal/domain/repo"
	"go-tpl/pkg/util"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type OtherRuleChecker struct {
	order         *dto.OrderDTO
	rulePriceRepo *repo.BiddingRulePriceRepo
}

func NewOtherRuleChecker(order *dto.OrderDTO, rulePriceRepo *repo.BiddingRulePriceRepo) *OtherRuleChecker {
	return &OtherRuleChecker{
		order:         order,
		rulePriceRepo: rulePriceRepo,
	}
}

func (c *OtherRuleChecker) CheckRule(ctx *gin.Context, ruleDTO *dto.RuleDTO) (*dto.BiddingResultDTO, error) {
	//ruleStr, _ := jsoniter.MarshalToString(rule)
	//logx.Infof(ctx, "check other rule: %s", ruleStr)

	rule := ruleDTO.OtherRule
	switch rule.RuleType {
	case 1: // 连锁模式
		return c.checkChain(ctx, ruleDTO)
	case 2: // 杂牌模式
		return c.checkSundry(ctx, ruleDTO)
	case 3: // 电影模式
		return c.checkMovie(ctx, ruleDTO)
	default:
		return nil, errors.New("未知规则类型")
	}
}

func (c *OtherRuleChecker) unifiedCheck(ctx *gin.Context, rule *entity.BiddingOtherRule) *dto.CheckFailCollector {
	failCollector := &dto.CheckFailCollector{}

	if rule.Status != 1 {
		failCollector.Add("报价规则未开启", "ruleStatus")
		return failCollector
	}

	if rule.LoverSeat != 0 && c.order.LoverSeat == 1 {
		failCollector.Add("情侣座不匹配", "loverSeat")
		return failCollector
	}
	return failCollector
}

func (c *OtherRuleChecker) checkChain(ctx *gin.Context, ruleDTO *dto.RuleDTO) (*dto.BiddingResultDTO, error) {
	rule := ruleDTO.OtherRule
	failCollector := c.unifiedCheck(ctx, rule)
	if failCollector.HasErr() {
		return nil, failCollector.GetErr()
	}

	if rule.CityName != "" {
		cityList := strings.Split(rule.CityName, ",")
		if !slices.Contains(cityList, c.order.CityName) {
			failCollector.Add("城市不匹配", "city")
			return nil, failCollector.GetErr()
		}
	}
	if rule.CityNotIn != "" {
		cityList := strings.Split(rule.CityNotIn, ",")
		if slices.Contains(cityList, c.order.CityName) {
			failCollector.Add("城市（不包含）不匹配", "city")
			return nil, failCollector.GetErr()
		}
	}

	// 判断影片类型
	if rule.MovieType != "" && c.order.MovieType != "" {
		if !strings.Contains(c.order.MovieType, rule.MovieType) {
			failCollector.Add("影片类型不匹配", "movieType")
			return nil, failCollector.GetErr()
		}
	}

	// 判断影院品牌
	if rule.CinemaBrandId != "" {
		cinemaIds := getCinemaIdsByBrandIds(util.SplitToInt64Slice(rule.CinemaBrandId, ","))
		if !util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("品牌不匹配", "brand")
			return nil, failCollector.GetErr()
		}
	}

	// 判断影院分类（包含、不包含）
	if rule.CinemaCategoryId != "" {
		cinemaIds := getCinemaIdsByCategoryIds(util.SplitToInt64Slice(rule.CinemaCategoryId, ","))
		if !util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("影院分类不匹配", "categoryId")
			return nil, failCollector.GetErr()
		}
	}
	if rule.NoCinemaCategoryId != "" {
		cinemaIds := getCinemaIdsByCategoryIds(util.SplitToInt64Slice(rule.NoCinemaCategoryId, ","))
		if util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("影院分类（不包含）不匹配", "categoryId")
			return nil, failCollector.GetErr()
		}
	}

	finalOffer := decimal.Zero

	priceList := ruleDTO.OtherRulePrice
	//priceList, _ := c.rulePriceRepo.ListByRuleId(ctx, rule.Id)
	if len(priceList) > 0 {
		for _, price := range priceList {
			if !checkHallClassify(c.order.HallClassifyIds, cast.ToString(price.HallId), "") {
				continue
			}
			if price.Fixed.GreaterThan(decimal.Zero) {
				finalOffer = price.Fixed
			}
		}
	}

	if rule.Discount.GreaterThan(decimal.Zero) {
		finalOffer = c.order.Price.Mul(rule.Discount.Div(decimal.NewFromInt(100))).Round(2)
	}

	if finalOffer.IsZero() {
		failCollector.Add("价格不匹配", "price")
		return nil, failCollector.GetErr()
	}

	var result = dto.BiddingResultDTO{
		Type:          "chain",
		UID:           rule.UserId,
		RuleId:        rule.Id,
		RuleName:      rule.Name,
		OtherRuleInfo: rule,
	}

	failReason, failType, _, _ := failCollector.Get()
	if len(failReason) == 0 && finalOffer.GreaterThan(decimal.Zero) {
		result.Status = 1
		result.Price = finalOffer
		result.IfBottomCover = rule.IfBottomCover == 1
	} else {
		result.Status = 0
		result.Price = decimal.Zero
		result.AllFailType = failType
		result.AllFailReason = failReason
	}
	return &result, nil
}

func (c *OtherRuleChecker) checkSundry(ctx *gin.Context, ruleDTO *dto.RuleDTO) (*dto.BiddingResultDTO, error) {
	var (
		rule      = ruleDTO.OtherRule
		priceList = ruleDTO.OtherRulePrice
	)
	failCollector := c.unifiedCheck(ctx, rule)
	if failCollector.HasErr() {
		return nil, failCollector.GetErr()
	}

	if rule.CinemaIds != "" {
		cinemaIds := util.SplitToInt64Slice(rule.CinemaIds, ",")
		if util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("影院（不包含）不匹配", "cinema")
			return nil, failCollector.GetErr()
		}
	}
	if rule.CinemaIdsIn != "" {
		cinemaIds := util.SplitToInt64Slice(rule.CinemaIdsIn, ",")
		if !util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("影院不匹配", "cinema")
			return nil, failCollector.GetErr()
		}
	}

	// 判断影院品牌
	if rule.CinemaBrandId != "" {
		cinemaIds := getCinemaIdsByBrandIds(util.SplitToInt64Slice(rule.CinemaBrandId, ","))
		if !util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("品牌不匹配", "brand")
			return nil, failCollector.GetErr()
		}
	}

	// 判断影院分类（包含、不包含）
	if rule.CinemaCategoryId != "" {
		cinemaIds := getCinemaIdsByCategoryIds(util.SplitToInt64Slice(rule.CinemaCategoryId, ","))
		if !util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("影院分类不匹配", "categoryId")
			return nil, failCollector.GetErr()
		}
	}
	if rule.NoCinemaCategoryId != "" {
		cinemaIds := getCinemaIdsByCategoryIds(util.SplitToInt64Slice(rule.NoCinemaCategoryId, ","))
		if util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("影院分类（不包含）不匹配", "categoryId")
			return nil, failCollector.GetErr()
		}
	}

	if !checkHallClassify(c.order.HallClassifyIds, rule.HallIds, "") {
		failCollector.Add("影厅不匹配", "hall")
		return nil, failCollector.GetErr()
	}

	var price *entity.BiddingRulePrice
	if len(priceList) > 0 {
		for _, rulePrice := range priceList {
			if rulePrice.Full.LessThanOrEqual(c.order.Price) {
				price = rulePrice
			}
		}
	}
	//price, _ := c.rulePriceRepo.GetByRuleIdAndPrice(ctx, rule.Id, c.order.Price)

	finalOffer := decimal.Zero
	if price != nil && price.Reduce.GreaterThan(decimal.Zero) {
		finalOffer = c.order.Price.Sub(price.Reduce)
	}
	if finalOffer.LessThanOrEqual(decimal.Zero) {
		if rule.Discount.IsZero() {
			failCollector.Add("价格不匹配", "price")
			return nil, failCollector.GetErr()
		} else {
			finalOffer = c.order.Price.Mul(rule.Discount.Div(decimal.NewFromInt(100))).Round(2)
		}
	}

	var result = dto.BiddingResultDTO{
		Type:          "sundry",
		UID:           rule.UserId,
		RuleId:        rule.Id,
		RuleName:      rule.Name,
		OtherRuleInfo: rule,
	}
	failReason, failType, _, _ := failCollector.Get()
	if len(failReason) == 0 && finalOffer.GreaterThan(decimal.Zero) {
		result.Status = 1
		result.Price = finalOffer
		result.IfBottomCover = rule.IfBottomCover == 1
	} else {
		result.Status = 0
		result.Price = decimal.Zero
		result.AllFailType = failType
		result.AllFailReason = failReason
	}
	return &result, nil
}

func (c *OtherRuleChecker) checkMovie(ctx *gin.Context, ruleDTO *dto.RuleDTO) (*dto.BiddingResultDTO, error) {
	var (
		rule      = ruleDTO.OtherRule
		priceList = ruleDTO.OtherRulePrice
	)

	failCollector := c.unifiedCheck(ctx, rule)
	if failCollector.HasErr() {
		return nil, failCollector.GetErr()
	}

	if !checkHallClassify(c.order.HallClassifyIds, rule.HallIds, "") {
		failCollector.Add("影厅不匹配", "hall")
		return nil, failCollector.GetErr()
	}

	if rule.MovieName != "" {
		if !strings.Contains(c.order.MovieName, rule.MovieName) {
			failCollector.Add("电影名称不匹配", "movieName")
			return nil, failCollector.GetErr()
		}
	}

	// 判断影院分类（包含、不包含）
	if rule.CinemaCategoryId != "" {
		cinemaIds := getCinemaIdsByCategoryIds(util.SplitToInt64Slice(rule.CinemaCategoryId, ","))
		if !util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("影院分类不匹配", "categoryId")
			return nil, failCollector.GetErr()
		}
	}
	if rule.NoCinemaCategoryId != "" {
		cinemaIds := getCinemaIdsByCategoryIds(util.SplitToInt64Slice(rule.NoCinemaCategoryId, ","))
		if util.Contains(cinemaIds, c.order.CinemaId) {
			failCollector.Add("影院分类（不包含）不匹配", "categoryId")
			return nil, failCollector.GetErr()
		}
	}

	//price, _ := c.rulePriceRepo.GetByRuleIdAndPrice2(ctx, rule.Id, c.order.Price, c.order.Price)
	var price *entity.BiddingRulePrice
	if len(priceList) > 0 {
		price = priceList[0]
		for _, rulePrice := range priceList {
			if rulePrice.Full.GreaterThan(c.order.Price) && rulePrice.Reduce.LessThanOrEqual(c.order.Price) {
				price = rulePrice
			}
		}
	}

	finalOffer := decimal.Zero
	if price != nil {
		finalOffer = price.Fixed
	}
	if finalOffer.LessThanOrEqual(decimal.Zero) {
		failCollector.Add("价格不匹配", "price")
		return nil, failCollector.GetErr()
	}

	var result = dto.BiddingResultDTO{
		Type:          "movie",
		UID:           rule.UserId,
		RuleId:        rule.Id,
		RuleName:      rule.Name,
		OtherRuleInfo: rule,
	}
	failReason, failType, _, _ := failCollector.Get()
	if len(failReason) == 0 && finalOffer.GreaterThan(decimal.Zero) {
		result.Status = 1
		result.Price = finalOffer
		result.IfBottomCover = rule.IfBottomCover == 1
	} else {
		result.Status = 0
		result.Price = decimal.Zero
		result.AllFailType = failType
		result.AllFailReason = failReason
	}
	return &result, nil
}

func checkHallClassify(orderHallIds, hallIds, noHallIds string) bool {
	orderHallIdList := util.SplitToInt64Slice(orderHallIds, ",")

	//没有匹配上普通厅再添加一个特殊厅匹配
	if len(orderHallIdList) > 0 && orderHallIdList[0] != 7 {
		orderHallIdList = append(orderHallIdList, 8)
	}

	containsOk, noContainsOk := true, true
	if hallIds != "" {
		hasHallIds := util.SplitToInt64Slice(hallIds, ",")
		containsOk = util.HasIntersection(hasHallIds, orderHallIdList)
	}

	if noHallIds != "" {
		noHasHallIds := util.SplitToInt64Slice(noHallIds, ",")
		noContainsOk = !util.HasIntersection(noHasHallIds, orderHallIdList)
	}
	return containsOk && noContainsOk
}
