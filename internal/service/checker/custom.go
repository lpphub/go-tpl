package checker

import (
	"go-tpl/internal/common/dto"
	"go-tpl/internal/domain/entity"
	"go-tpl/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/lpphub/golib/logger/logx"
	"github.com/shopspring/decimal"
)

type CustomRuleChecker struct {
	order *dto.OrderDTO
}

func NewCustomRuleChecker(order *dto.OrderDTO) *CustomRuleChecker {
	return &CustomRuleChecker{
		order: order,
	}
}

func (c *CustomRuleChecker) CheckRule(ctx *gin.Context, rule *entity.BiddingCustomRule) (*dto.BiddingResultDTO, error) {
	ruleStr, _ := jsoniter.MarshalToString(rule)
	logx.Infof(ctx, "check custom rule: %s", ruleStr)

	failCollector := &dto.CheckFailCollector{}

	resCity := c.checkCity(ctx, rule)
	if !resCity.Passed {
		failCollector.Add(resCity.FailReason, resCity.FailType)
		return nil, failCollector.GetErr()
	}

	resCinema := c.checkCinema(ctx, rule)
	if !resCinema.Passed {
		failCollector.Add(resCinema.FailReason, resCinema.FailType)
		return nil, failCollector.GetErr()
	}

	resHall := c.checkHall(ctx, rule)
	if !resHall.Passed {
		failCollector.Add(resHall.FailReason, resHall.FailType)
		return nil, failCollector.GetErr()
	}

	resMovie := c.checkMovie(ctx, rule)
	if !resMovie.Passed {
		failCollector.Add(resMovie.FailReason, resMovie.FailType)
		return nil, failCollector.GetErr()
	}

	resMovieType := c.checkMovieType(ctx, rule)
	if !resMovieType.Passed {
		failCollector.Add(resMovieType.FailReason, resMovieType.FailType)
		return nil, failCollector.GetErr()
	}

	resPrice := c.checkPrice(ctx, rule)
	if !resPrice.Passed {
		failCollector.Add(resPrice.FailReason, resPrice.FailType)
		return nil, failCollector.GetErr()
	}

	resSeatNum := c.checkSeatNum(ctx, rule)
	if !resSeatNum.Passed {
		failCollector.Add(resSeatNum.FailReason, resSeatNum.FailType)
		return nil, failCollector.GetErr()
	}

	resLoverSeat := c.checkLoverSeat(ctx, rule)
	if !resLoverSeat.Passed {
		failCollector.Add(resLoverSeat.FailReason, resLoverSeat.FailType)
		return nil, failCollector.GetErr()
	}

	resWeekNum := c.checkWeekNum(ctx, rule)
	if !resWeekNum.Passed {
		failCollector.Add(resWeekNum.FailReason, resWeekNum.FailType)
		return nil, failCollector.GetErr()
	}

	resHourSpan := c.checkShowTime(ctx, rule)
	if !resHourSpan.Passed {
		failCollector.Add(resHourSpan.FailReason, resHourSpan.FailType)
		return nil, failCollector.GetErr()
	}

	resCinemaCat := c.checkCinemaCategory(ctx, rule)
	if !resCinemaCat.Passed {
		failCollector.Add(resCinemaCat.FailReason, resCinemaCat.FailType)
		return nil, failCollector.GetErr()
	}

	resBrand := c.checkBrand(ctx, rule)
	if !resBrand.Passed {
		failCollector.Add(resBrand.FailReason, resBrand.FailType)
		return nil, failCollector.GetErr()
	}

	resHallClassify := c.checkHallClassify(ctx, rule)
	if !resHallClassify.Passed {
		failCollector.Add(resHallClassify.FailReason, resHallClassify.FailType)
		return nil, failCollector.GetErr()
	}

	resIsSeat := c.checkIsSeat(ctx, rule)
	if !resIsSeat.Passed {
		failCollector.Add(resIsSeat.FailReason, resIsSeat.FailType)
		return nil, failCollector.GetErr()
	}

	allFailReasons, allFailTypes, lastFailReason, _ := failCollector.Get()
	allPassed := len(allFailReasons) == 0 // 无失败即全部通过

	var result = dto.BiddingResultDTO{
		Type:     "customize",
		UID:      rule.Uid,
		RuleId:   rule.Id,
		RuleName: rule.Name,
		RuleInfo: rule,
	}
	if allPassed {
		price := c.calculatePrice(rule)
		result.Status = 1
		result.Price = price
		result.IfBottomCover = rule.IfBottomCover == 1
	} else {
		result.Status = 0
		result.Price = decimal.Zero
		result.AllFailReason = allFailReasons
		result.AllFailType = allFailTypes
		result.FailReason = lastFailReason
	}
	return &result, nil
}

func (c *CustomRuleChecker) checkCity(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.CityName == "" {
		return dto.CheckSubResultDTO{Passed: true} // 未配置则通过
	}
	if rule.CityName != "" {
		cityList := util.SplitNonEmpty(rule.CityName, ",")
		if !util.Contains(cityList, c.order.CityName) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "城市不匹配",
				FailType:   "city",
			}
		}
	}

	if rule.NoCityName != "" {
		cityList := util.SplitNonEmpty(rule.NoCityName, ",")
		if util.Contains(cityList, c.order.CityName) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "城市（不包含）不匹配",
				FailType:   "city",
			}
		}
	}
	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkCinema(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.CinemaId != "" {
		cinemaList := util.SplitToInt64Slice(rule.CinemaId, ",")
		if !util.Contains(cinemaList, c.order.CinemaId) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "影院不匹配",
				FailType:   "cinema",
			}
		}
	}

	if rule.NoCinemaId != "" {
		cinemaList := util.SplitToInt64Slice(rule.NoCinemaId, ",")
		if util.Contains(cinemaList, c.order.CinemaId) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "影院（不包含）不匹配",
				FailType:   "cinema",
			}
		}
	}
	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkHall(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.HallName != "" {
		if !util.SplitAndMatch(rule.HallName, c.order.HallName) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "影厅不匹配",
				FailType:   "hall",
			}
		}
	}

	if rule.NoHallName != "" {
		if util.SplitAndMatch(rule.NoHallName, c.order.HallName) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "影厅（不包含）不匹配",
				FailType:   "hall",
			}
		}
	}
	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkMovie(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.MovieName != "" {
		movieList := util.SplitNonEmpty(rule.MovieName, ",")
		if !util.Contains(movieList, c.order.MovieName) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "电影不匹配",
				FailType:   "movie",
			}
		}
	}

	if rule.NoMovieName != "" {
		movieList := util.SplitNonEmpty(rule.NoMovieName, ",")
		if util.Contains(movieList, c.order.MovieName) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "电影（不包含）不匹配",
				FailType:   "movie",
			}
		}
	}
	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkPrice(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	var passed bool

	if rule.Type == 4 { // 总市场价格区间
		totalPrice := c.order.MaoyanPrice.Mul(decimal.NewFromInt(int64(c.order.SeatNum)))
		passed = rule.PriceMin.GreaterThan(decimal.Zero) && rule.PriceMin.LessThanOrEqual(totalPrice)
	} else { // 单张市场价格区间
		if rule.PriceMin.GreaterThan(decimal.Zero) || rule.PriceMax.GreaterThan(decimal.Zero) {
			minPrice := c.order.MaoyanPrice
			maxPrice := c.order.MaoyanPrice

			switch {
			case rule.PriceMin.GreaterThan(decimal.Zero) && rule.PriceMax.Equal(decimal.Zero) && rule.PriceMin.LessThanOrEqual(minPrice):
				passed = true
			case rule.PriceMax.GreaterThan(decimal.Zero) && rule.PriceMin.Equal(decimal.Zero) && rule.PriceMax.GreaterThanOrEqual(maxPrice):
				passed = true
			case rule.PriceMax.GreaterThan(decimal.Zero) && rule.PriceMin.GreaterThan(decimal.Zero) &&
				rule.PriceMax.GreaterThanOrEqual(maxPrice) && rule.PriceMin.LessThanOrEqual(minPrice):
				passed = true
			default:
				passed = false
			}
		} else {
			passed = true
		}
	}

	if !passed {
		return dto.CheckSubResultDTO{
			Passed:     false,
			FailReason: "价格不匹配",
			FailType:   "price",
		}
	}

	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkSeatNum(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.Num != "" && rule.Num != "0" {
		seatNumList := util.SplitToInt64Slice(rule.Num, ",")
		if !util.Contains(seatNumList, int64(c.order.SeatNum)) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "座位数不匹配",
				FailType:   "seatNum",
			}
		}
	}
	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkLoverSeat(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.LoverSeat != 1 || c.order.LoverSeat != 0 {
		return dto.CheckSubResultDTO{
			Passed:     false,
			FailReason: "情侣座不匹配",
			FailType:   "loverSeat",
		}
	}

	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkIsSeat(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.IsSeat != 1 || !util.Contains([]int8{1, 2, 3}, c.order.IsSeat) {
		return dto.CheckSubResultDTO{
			Passed:     false,
			FailReason: "可调座不匹配",
			FailType:   "isSeat",
		}
	}

	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkMovieType(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.MovieType != "" {
		if !util.ContainsNoCase(c.order.MovieType, rule.MovieType) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "影片类型不匹配",
				FailType:   "movieType",
			}
		}
	}

	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkWeekNum(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.WeekNum != "" {
		weekNum := int64(time.Unix(int64(c.order.ShowTime), 0).Weekday())
		if weekNum == 0 { // 周日转换为7
			weekNum = 7
		}

		weekNumList := util.SplitToInt64Slice(rule.WeekNum, ",")
		if !util.Contains(weekNumList, weekNum) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "包含星期不匹配",
				FailType:   "weekNum",
			}
		}
	}
	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkShowTime(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.HourSpan != "" {
		spanParts := util.SplitNonEmpty(rule.HourSpan, "-")
		if len(spanParts) != 2 {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "放映时间配置无效",
				FailType:   "hourSpan",
			}
		}

		parts := util.SplitNonEmpty(rule.HourSpan, "-")
		if len(parts) != 2 {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "放映时间配置无效",
				FailType:   "hourSpan",
			}
		}
		st, serr := time.Parse("15:04", parts[0])
		et, eerr := time.Parse("15:04", parts[1])
		if serr != nil || eerr != nil {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "放映时间解析异常",
				FailType:   "hourSpan",
			}
		}

		movieShowTime := time.Unix(int64(c.order.ShowTime), 0)
		if movieShowTime.Before(st) || movieShowTime.After(et) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "放映时间不匹配",
				FailType:   "hourSpan",
			}
		}

	}
	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkBrand(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	if rule.BrandId != "" && rule.BrandId != "0" {
		brandIds := util.SplitToInt64Slice(rule.BrandId, ",")
		if len(brandIds) == 0 {
			return dto.CheckSubResultDTO{Passed: true}
		}

		cinemaIds := getCinemaIdsByBrandIds(brandIds)
		if !util.Contains(cinemaIds, c.order.CinemaId) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "品牌不匹配",
				FailType:   "brand",
			}
		}
	}

	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkHallClassify(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	ok := checkHallClassify(c.order.HallClassifyIds, rule.HallIds, rule.NoHallIds)
	if !ok {
		return dto.CheckSubResultDTO{
			Passed:     false,
			FailReason: "影厅类型不匹配",
			FailType:   "hall",
		}
	}
	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) checkCinemaCategory(_ *gin.Context, rule *entity.BiddingCustomRule) dto.CheckSubResultDTO {
	orderCinemaId := c.order.CinemaId
	// 检查包含影院分类
	if rule.CinemaCategoryId != "" {
		cinemaIds := getCinemaIdsByCategoryIds(util.SplitToInt64Slice(rule.CinemaCategoryId, ","))
		if !util.Contains(cinemaIds, orderCinemaId) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "影院分类不匹配",
				FailType:   "categoryId",
			}
		}
	}

	// 检查不包含影院分类规则
	if rule.NoCinemaCategoryId != "" {
		cinemaIds := getCinemaIdsByCategoryIds(util.SplitToInt64Slice(rule.NoCinemaCategoryId, ","))
		if !util.Contains(cinemaIds, orderCinemaId) {
			return dto.CheckSubResultDTO{
				Passed:     false,
				FailReason: "影院分类（不包含）不匹配",
				FailType:   "categoryId",
			}
		}
	}

	return dto.CheckSubResultDTO{Passed: true}
}

func (c *CustomRuleChecker) calculatePrice(rule *entity.BiddingCustomRule) decimal.Decimal {
	var price decimal.Decimal

	switch rule.Type {
	case 0:
		// 0=市场价比例报价: 单价 × (比例值/100)
		ratio := rule.Value.Div(decimal.NewFromInt(100))
		price = c.order.MaoyanPrice.Mul(ratio)
	case 1:
		// 1=最高价减法: (单价 × 最高价系数) - 固定值
		price = c.order.MaoyanPrice.Mul(c.order.Max).Sub(rule.Value)
	case 2:
		// 2=市场价减法: 单价 - 固定值
		price = c.order.MaoyanPrice.Sub(rule.Value)
	case 4:
		// 4=市场总价满减: (总价 - 固定值) ÷ 座位数
		totalPrice := c.order.MaoyanPrice.Mul(decimal.NewFromInt(int64(c.order.SeatNum))).Sub(rule.Value)
		price = totalPrice.Div(decimal.NewFromInt(int64(c.order.SeatNum)))
	default:
		// 其他=固定金额
		price = rule.Value
	}
	rounded := price.Mul(decimal.NewFromInt(10)).Ceil().Div(decimal.NewFromInt(10))
	return rounded
}
