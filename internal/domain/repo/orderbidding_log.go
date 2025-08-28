package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderBiddingLogRepo struct {
	db *gorm.DB
}

func NewOrderBiddingLogRepo(db *gorm.DB) *OrderBiddingLogRepo {
	return &OrderBiddingLogRepo{db: db}
}

func (r *OrderBiddingLogRepo) Create(ctx *gin.Context, log *entity.OrderAutoBiddingLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
