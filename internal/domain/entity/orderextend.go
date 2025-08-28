package entity

import "github.com/shopspring/decimal"

type OrderExtend struct {
	Id                      int64           `gorm:"column:id;primaryKey"`             // 主键ID（数据库为 bigint unsigned，对应Go的uint64）
	OrderId                 int64           `gorm:"column:orderId"`                   // 关联订单ID（对应 tb_order_jj.id）
	OrderNo                 string          `gorm:"column:order_no"`                  // 订单编号（关联订单的订单号）
	CustomerPayType         string          `gorm:"column:customer_pay_type"`         // 个人购票订单支付类型
	PayMoney                decimal.Decimal `gorm:"column:pay_money"`                 // 客户实际支付金额
	RefundMoney             decimal.Decimal `gorm:"column:refund_money"`              // 退款金额
	HallClassifyId          int64           `gorm:"column:hall_classify_id"`          // 影厅ID（单个）
	HallClassifyIds         string          `gorm:"column:hall_classify_ids"`         // 影厅ID（多个匹配模式，字符串存储）
	CreateTime              string          `gorm:"column:createTime"`                // 创建时间（datetime格式，对应Go的string）
	UpdateTime              string          `gorm:"column:updateTime"`                // 更新时间（自动更新datetime，对应Go的string）
	RetailerAccountId       int64           `gorm:"column:retailer_account_id"`       // 分销商账号ID
	CouponCode              string          `gorm:"column:coupon_code"`               // 优惠券编码
	CouponMoney             decimal.Decimal `gorm:"column:coupon_money"`              // 优惠券抵扣金额
	RetailerCommission      decimal.Decimal `gorm:"column:retailer_commission"`       // 分销商佣金
	TotalRetailerCommission decimal.Decimal `gorm:"column:total_retailer_commission"` // 总分销佣金（分销商+二级分销）
	TicketMethod            int8            `gorm:"column:ticket_method"`             // 出票方式：1.竞价出票，2.自动出票
	ShowTicketCount         int             `gorm:"column:show_ticket_count"`         // 查看取票码次数
	SendTicketStatus        int8            `gorm:"column:send_ticket_status"`        // 取票码发送状态：1.已发送，2.未发送
	AutoAssignStatus        int8            `gorm:"column:auto_assign_status"`        // 自动分配状态：1.开启，2.关闭
	RealTotalSellPrice      decimal.Decimal `gorm:"column:real_total_sell_price"`     // 真实总售价
	LockEndTime             int             `gorm:"column:lock_end_time"`             // 锁座自动解锁时间（时间戳）
	LockPlatform            string          `gorm:"column:lock_platform"`             // 锁座平台（如特定影院系统标识）
	ChannelOrderNo          string          `gorm:"column:channel_order_no"`          // 出票渠道订单号
	SeatRegionName          string          `gorm:"column:seat_region_name"`          // 座位区域名称（如"前排"、"VIP区"）
	SeatRegionId            string          `gorm:"column:seat_region_id"`            // 座位区域ID
	MovieType               string          `gorm:"column:movie_type"`                // 电影制式（如"2D"、"3D"、"IMAX"）
}

func (OrderExtend) TableName() string {
	return "tb_order_jj_extend"
}
