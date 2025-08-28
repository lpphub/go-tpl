package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BiddingCustomRuleRepo struct {
	db *gorm.DB
}

func NewBiddingCustomRuleRepo(db *gorm.DB) *BiddingCustomRuleRepo {
	return &BiddingCustomRuleRepo{
		db: db,
	}
}

// ListWithEffective 获取有效状态的规则
func (r *BiddingCustomRuleRepo) ListWithEffective(ctx *gin.Context) ([]*entity.BiddingCustomRule, error) {
	var list []*entity.BiddingCustomRule
	err := r.db.WithContext(ctx).Model(&entity.BiddingCustomRule{}).
		Where("id = 328710 and status =? and delete_time =?", 1, 0).
		Order("id desc").
		Find(&list).Error
	return list, err
}
