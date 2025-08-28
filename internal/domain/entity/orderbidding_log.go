package entity

import "go-tpl/pkg/ext"

type OrderAutoBiddingLog struct {
	Id         int64          `gorm:"column:id;primaryKey"` // 主键ID（数据库为 int unsigned，对应Go的uint64）
	OrderId    int64          `gorm:"column:orderId"`       // 关联订单ID（追溯日志对应的订单）
	Uid        int64          `gorm:"column:uid"`           // 票商ID（关联用户表，记录操作票商）
	RuleId     int64          `gorm:"column:ruleId"`        // 关联自动竞价规则ID（追溯日志对应的竞价规则）
	RuleName   string         `gorm:"column:ruleName"`      // 规则名称（日志备份用，与原规则名称一致）
	RuleType   string         `gorm:"column:ruleType"`      // 规则类型（如“连锁模式”“电影模式”，原规则配置备份）
	RuleInfo   string         `gorm:"column:ruleInfo"`      // 规则详情（文本存储完整规则配置，便于审计追溯）
	CreateTime *ext.Timestamp `gorm:"column:createTime"`    // 日志创建时间（自动生成datetime，记录竞价操作时间）
	UpdateTime *ext.Timestamp `gorm:"column:updateTime"`    // 日志更新时间（自动更新datetime，冗余字段）
}

func (OrderAutoBiddingLog) TableName() string {
	return "tb_hh_order_auto_bidding_log"
}
