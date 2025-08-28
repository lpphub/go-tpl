package entity

import "github.com/shopspring/decimal"

type BiddingRulePrice struct {
	Id       int64           `gorm:"column:id;primaryKey"` // 主键ID（自增）
	UserId   int64           `gorm:"column:user_id"`       // 关联用户ID
	RuleId   int64           `gorm:"column:rule_id"`       // 关联竞价规则ID（对应 tb_hh_auto_bidding_rule.id，核心关联字段）
	RuleType int             `gorm:"column:rule_type"`     // 规则模式：1.连锁模式，2.杂牌模式，3.电影模式（与关联规则的模式一致）
	HallName string          `gorm:"column:hall_name"`     // 影厅名称
	HallId   int64           `gorm:"column:hall_id"`       // 影厅分类ID
	Fixed    decimal.Decimal `gorm:"column:fixed"`         // 固定报价金额
	Full     decimal.Decimal `gorm:"column:full"`          // 满减门槛金额
	Reduce   decimal.Decimal `gorm:"column:reduce"`        // 满减金额
}

func (BiddingRulePrice) TableName() string {
	return "tb_hh_auto_bidding_rule_price"
}
