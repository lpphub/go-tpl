package repo

import (
	"go-tpl/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BiddingRulePriceRepo struct {
	db *gorm.DB
}

func NewBiddingRulePriceRepo(db *gorm.DB) *BiddingRulePriceRepo {
	return &BiddingRulePriceRepo{db: db}
}

func (r *BiddingRulePriceRepo) GetByRuleIdAndPrice(ctx *gin.Context, ruleId int64, orderPrice decimal.Decimal) (*entity.BiddingRulePrice, error) {
	var price *entity.BiddingRulePrice
	err := r.db.WithContext(ctx).Where("rule_id = ? and full <= ?", ruleId, orderPrice).Order("full desc").First(&price).Error
	return price, err
}

func (r *BiddingRulePriceRepo) ListByRuleId(ctx *gin.Context, ruleId int64) ([]*entity.BiddingRulePrice, error) {
	var prices []*entity.BiddingRulePrice
	err := r.db.WithContext(ctx).Where("rule_id = ?", ruleId).Order("full desc").Find(&prices).Error
	return prices, err
}

func (r *BiddingRulePriceRepo) GetByRuleIdAndPrice2(ctx *gin.Context, ruleId int64, orderPrice, reduce decimal.Decimal) (*entity.BiddingRulePrice, error) {
	var price *entity.BiddingRulePrice
	err := r.db.WithContext(ctx).Where("rule_id = ? and full > ? and reduce <= ?", ruleId, orderPrice, reduce).Order("fixed").First(&price).Error
	return price, err
}

func (r *BiddingRulePriceRepo) ListByRuleIds(ctx *gin.Context, ruleIds []int64) ([]*entity.BiddingRulePrice, error) {
	var prices []*entity.BiddingRulePrice
	err := r.db.WithContext(ctx).Where("rule_id in ?", ruleIds).Order("full desc").Find(&prices).Error
	return prices, err
}
