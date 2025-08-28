package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderExtendRepo struct {
	db *gorm.DB
}

func NewOrderExtendRepo(db *gorm.DB) *OrderExtendRepo {
	return &OrderExtendRepo{db: db}
}

func (o *OrderExtendRepo) GetOne(ctx *gin.Context, orderId int64) (*entity.OrderExtend, error) {
	var orderExtend entity.OrderExtend
	err := o.db.WithContext(ctx).Where("orderId = ?", orderId).First(&orderExtend).Error
	if err != nil {
		return nil, err
	}
	return &orderExtend, nil
}
