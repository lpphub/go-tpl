package entity

import "go-tpl/pkg/ext"

type BiddingBillBottomCover struct {
	Id         int64          `gorm:"column:id;primaryKey"` // 主键ID
	BiddingId  int64          `gorm:"column:bidding_id"`    // 关联竞价ID
	OrderId    int64          `gorm:"column:order_id"`      // 关联订单ID
	Uid        int64          `gorm:"column:uid"`           // 票商ID
	CreateTime *ext.Timestamp `gorm:"column:create_time"`   // 记录创建时间
	UpdateTime *ext.Timestamp `gorm:"column:update_time"`   // 记录更新时间
}

func (BiddingBillBottomCover) TableName() string {
	return "tb_bidding_bill_bottom_cover"
}
