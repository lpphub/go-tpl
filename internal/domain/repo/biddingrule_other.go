package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BiddingOtherRuleRepo struct {
	db *gorm.DB
}

func NewBiddingOtherRuleRepo(db *gorm.DB) *BiddingOtherRuleRepo {
	return &BiddingOtherRuleRepo{db: db}
}

func (r *BiddingOtherRuleRepo) ListWithEffective(ctx *gin.Context) ([]*entity.BiddingOtherRule, error) {
	var list []*entity.BiddingOtherRule
	err := r.db.WithContext(ctx).Where("status = ?", 1).Find(&list).Error
	return list, err
}
