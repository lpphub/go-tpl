package checker

import (
	"go-tpl/internal/common/constant"
	"go-tpl/internal/common/dto"
	"go-tpl/internal/domain/entity"
	"go-tpl/internal/domain/repo"
	"go-tpl/pkg/util"

	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

var levelBonus = map[int]int{
	1: 10,
	2: 20,
	3: 30,
}

type UserChecker struct {
	order    *dto.OrderDTO
	userRepo *repo.UserRepo
	billRepo *repo.BiddingBillRepo
	redis    *redis.Client
}

func NewUserChecker(order *dto.OrderDTO, userRepo *repo.UserRepo, billRepo *repo.BiddingBillRepo, redis *redis.Client) *UserChecker {
	return &UserChecker{
		order:    order,
		userRepo: userRepo,
		billRepo: billRepo,
		redis:    redis,
	}
}

func (c *UserChecker) Check(ctx *gin.Context, ruleResultList []*dto.BiddingResultDTO, isBidding bool) ([]*dto.BiddingResultDTO, error) {
	var result []*dto.BiddingResultDTO

	userRuleMap := make(map[int64][]*dto.BiddingResultDTO)
	for _, rule := range ruleResultList {
		userRuleMap[rule.UID] = append(userRuleMap[rule.UID], rule)
	}

	// todo 优化：用户信息从redis缓存获取
	users, _ := c.userRepo.ListEffectiveByIds(ctx, util.MapKeyToSlice(userRuleMap))

	userJoinList := c.redis.SMembers(ctx, fmt.Sprintf(constant.CacheBiddingJoin, c.order.Id)).Val()
	alreadyJoinUser := make(map[int64]bool)
	if len(userJoinList) > 0 {
		for _, uid := range userJoinList {
			id := cast.ToInt64(uid)
			alreadyJoinUser[id] = true
		}
	}
	for i, user := range users {
		if user.UserId == 81654 && c.order.EndTime <= time.Now().Unix() {
			logx.Infof(ctx, "api票商不报价: %d", user.UserId)
			break
		}

		if err := c.checkLimitTime(ctx, user); err != nil {
			logx.Error(ctx, err.Error())
			break
		}
		if isBidding {
			if alreadyJoinUser[user.UserId] {
				logx.Infof(ctx, "用户已经报过价: %d", user.UserId)
				break
			}

			if err := c.checkBidTime(ctx, user); err != nil {
				logx.Error(ctx, err.Error())
				break
			}

			if err := c.checkBidMaxTimes(ctx, user); err != nil {
				logx.Error(ctx, err.Error())
				break
			}
		}

		// 用户规则检验没问题，将其对应的规则添加进结果
		for _, rule := range userRuleMap[user.UserId] {
			rule.UserInfo = users[i]
			if !rule.IfBottomCover {
				rule.IfBottomCover = user.BottomCover == 1
			}
			result = append(result, rule)
		}
	}

	// todo 增加黑名单校验
	return result, nil
}

func (c *UserChecker) checkLimitTime(_ *gin.Context, user *entity.User) error {
	if user.IsAutoBiddingLimitOrder == 0 && c.order.LimitTime > 0 && c.order.LimitTime-c.order.CreateTime <= 1800 {
		return errors.New("用户不能报价限时单")
	}
	return nil
}

func (c *UserChecker) checkAlreadyBid(ctx *gin.Context, user *entity.User) error {
	result, err := c.redis.SIsMember(ctx, fmt.Sprintf(constant.CacheBiddingJoin, c.order.Id), user.UserId).Result()
	if err != nil {
		return errors.WithMessage(err, "查询用户是否报过价失败")
	}
	if result {
		return errors.New("用户已经报过价")
	}
	return nil
}

func (c *UserChecker) checkBidTime(_ *gin.Context, user *entity.User) error {
	timeParts := strings.Split(user.BiddingTime, ",")

	now := time.Now()
	date := now.Format("2006-01-02")

	start, err := time.Parse("2006-01-02 15:04", date+" "+timeParts[0])
	if err != nil {
		return err
	}
	end, err := time.Parse("2006-01-02 15:04", date+" "+timeParts[1])
	if err != nil {
		return err
	}

	inRange := false
	if start.After(end) {
		inRange = now.After(start) || now.Before(end)
	} else {
		inRange = (now.After(start) || now.Equal(start)) &&
			(now.Before(end) || now.Equal(end))
	}
	if !inRange {
		return errors.New("用户该时段不能报价")
	}
	return nil
}

func (c *UserChecker) checkBidMaxTimes(ctx *gin.Context, user *entity.User) error {
	count := user.MaxBiddingTimes
	if count <= 0 {
		count = levelBonus[user.Level]
	}
	// todo 优化：避免查库
	billList, err := c.billRepo.ListByUidAndNotOrderId(ctx, user.UserId, c.order.Id)
	if err != nil {
		return err
	}

	userBiddingCount := make(map[int64]int)
	for _, item := range billList {
		if item.IsAuto == 0 {
			userBiddingCount[item.Uid]++
		} else {
			ifEffectiveBidding := false
			if item.Uid == 81654 { // 特殊用户（表示有效竞价）
				ifEffectiveBidding = true
			} else {
				ifEffectiveBidding = c.redis.SIsMember(ctx, fmt.Sprintf(constant.CacheBiddingEffectiveRanking, item.OrderId), item.Uid).Val()
			}

			if ifEffectiveBidding {
				userBiddingCount[item.Uid]++
			}
		}
	}

	if userBiddingCount[user.UserId] >= count {
		return errors.New("用户报价次数达到上限")
	}
	return nil
}
