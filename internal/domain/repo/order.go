package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (o *OrderRepo) GetOne(ctx *gin.Context, orderId int64) (*entity.Order, error) {
	var order entity.Order
	// todo 补状态
	err := o.db.WithContext(ctx).Where("id = ?", orderId).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
