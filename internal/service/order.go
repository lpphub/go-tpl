package service

import (
	"go-tpl/internal/api/vo"
	"go-tpl/internal/common/constant"
	"go-tpl/internal/common/dto"
	"go-tpl/internal/domain/repo"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

type OrderService struct {
	orderRepo  *repo.OrderRepo
	extendRepo *repo.OrderExtendRepo
	redis      *redis.Client
}

func NewOrderService(orderRepo *repo.OrderRepo, extendRepo *repo.OrderExtendRepo, redis *redis.Client) *OrderService {
	return &OrderService{
		orderRepo:  orderRepo,
		extendRepo: extendRepo,
		redis:      redis,
	}
}

func (s *OrderService) GetOneForBidding(ctx *gin.Context, params vo.BiddingAutoDTO) (*dto.OrderDTO, error) {
	order, err := s.orderRepo.GetOne(ctx, params.OrderId)
	if err != nil {
		return nil, err
	}
	extend, err := s.extendRepo.GetOne(ctx, params.OrderId)
	if err != nil {
		return nil, err
	}

	orderDto := &dto.OrderDTO{
		Order:           *order,
		HallClassifyId:  extend.HallClassifyId,
		HallClassifyIds: extend.HallClassifyIds,
		MovieType:       extend.MovieType,
		IsBottomCover:   params.IsOrderBottomCover,
	}
	if params.OrderLimitTime > 0 {
		order.LimitTime = params.OrderLimitTime
	}

	val := s.redis.ZScore(ctx, constant.CacheOrderEndTime, cast.ToString(order.Id)).Val()
	orderDto.EndTime = cast.ToInt64(val)
	return orderDto, nil
}
