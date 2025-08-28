package entity

import "github.com/shopspring/decimal"

type Order struct {
	Id                        int64           `gorm:"column:id;primaryKey"`                // 主键ID
	MerchantId                int64           `gorm:"column:merchant_id"`                  // 商家id
	UserId                    int64           `gorm:"column:user_id"`                      // 用户id
	OrderId                   string          `gorm:"column:order_id"`                     // 订单编号（唯一索引）
	OrderType                 string          `gorm:"column:order_type"`                   // 订单类型
	PayType                   int8            `gorm:"column:pay_type"`                     // 支付类型: 0券商 1竞价
	CinemaId                  int64           `gorm:"column:cinemaId"`                     // 影院id
	CityId                    int64           `gorm:"column:cityId"`                       // 城市id
	MovieInfo                 string          `gorm:"column:movieInfo"`                    // 电影信息（文本）
	PlanInfo                  string          `gorm:"column:planInfo"`                     // 场次信息（文本）
	Image                     string          `gorm:"column:image"`                        // 图片信息（文本）
	Seats                     string          `gorm:"column:seats"`                        // 座位信息（文本）
	InvalidateDate            string          `gorm:"column:invalidateDate"`               // 电影开场日期
	StartTime                 string          `gorm:"column:startTime"`                    // 场次时间
	TotalPrice                decimal.Decimal `gorm:"column:total_price"`                  // 总金额（淘宝总售价）
	Price                     decimal.Decimal `gorm:"column:price"`                        // 价格（淘宝售价）
	OriginalPrice             decimal.Decimal `gorm:"column:originalPrice"`                // 原价
	PayPrice                  decimal.Decimal `gorm:"column:payPrice"`                     // 总成本价（票商出价+平台抽成）
	CostPrice                 decimal.Decimal `gorm:"column:cost_price"`                   // 成本价（票商出价）
	MaoyanPrice               decimal.Decimal `gorm:"column:maoyan_price"`                 // 猫眼价格
	SeatNum                   int             `gorm:"column:seat_num"`                     // 座位数量
	SeatInfo                  string          `gorm:"column:seatInfo"`                     // 座位信息（长字符串，1500字符）
	Status                    int8            `gorm:"column:status"`                       // 订单状态: 0竞价中 1待出票(带上报) 2已出票(带结算) 3出票失败 4取消订单 5已结算 6纠纷中 7关闭
	FailType                  int8            `gorm:"column:fail_type"`                    // 失败类型：1.超时失败
	CreateTime                int             `gorm:"column:create_time"`                  // 创建时间（时间戳）
	Phone                     string          `gorm:"column:phone"`                        // 联系电话
	Tphone                    string          `gorm:"column:tphone"`                       // 用户手机号
	HallName                  string          `gorm:"column:hallName"`                     // 影厅名称（文本）
	Refund                    int8            `gorm:"column:refund"`                       // 退款状态: 0未退款 , 1已退款
	IsDel                     int8            `gorm:"column:is_del"`                       // 是否已删除
	CinemaName                string          `gorm:"column:cinemaName"`                   // 影院名称
	CityName                  string          `gorm:"column:cityName"`                     // 市级信息
	ProvinceName              string          `gorm:"column:provinceName"`                 // 省级信息
	MovieName                 string          `gorm:"column:movieName"`                    // 电影名称
	Address                   string          `gorm:"column:address"`                      // 影院地址
	BiddId                    int64           `gorm:"column:bidd_id"`                      // 绑定竞价ID
	ToType                    int8            `gorm:"column:toType"`                       // 购票类型:0实单,1虚单
	IsSeat                    int8            `gorm:"column:is_seat"`                      // 调座:0不可调，1可调，2九宫格，3协商调
	IsAuto                    int8            `gorm:"column:is_auto"`                      // 自动报价标识
	LoverSeat                 int8            `gorm:"column:loverSeat"`                    // 情侣座标识
	Paid                      int             `gorm:"column:paid"`                         // 支付标记(三方)
	PayTime                   int             `gorm:"column:pay_time"`                     // 支付时间（时间戳）
	Dwz                       string          `gorm:"column:dwz"`                          // 短网址
	ToCount                   int             `gorm:"column:to_count"`                     // 循环次数
	Type                      int8            `gorm:"column:type"`                         // 出票方式:0竞价,1自动
	AutoNo                    string          `gorm:"column:autoNo"`                       // 自动出票订单号（索引）
	Min                       decimal.Decimal `gorm:"column:min"`                          // 最小值（金额）
	Max                       decimal.Decimal `gorm:"column:max"`                          // 最大值（金额）
	Img                       string          `gorm:"column:img"`                          // 识别图片信息（文本）
	FormType                  string          `gorm:"column:form_type"`                    // cinemaId、cityId 来源
	ManId                     string          `gorm:"column:man_id"`                       // 趣满满订单ID
	ManPrice                  decimal.Decimal `gorm:"column:man_price"`                    // 趣满满扣价
	ServiceCost               string          `gorm:"column:servicecost"`                  // 趣满满服务成本
	ManMsg                    string          `gorm:"column:man_msg"`                      // 趣满满失败内容
	IsLock                    int8            `gorm:"column:is_lock"`                      // 是否锁座
	IsFrom                    int8            `gorm:"column:is_from"`                      // 录入信息来源: 0总后台,1PC,2微信端,3自录单,4客户自助下单,5商家单
	IsType                    int8            `gorm:"column:is_type"`                      // 真假失败（默认假，舍弃字段）
	IsErr                     int8            `gorm:"column:is_err"`                       // 是否处理
	IsAdmin                   int8            `gorm:"column:is_admin"`                     // 手动失败标识
	BidTime                   int             `gorm:"column:bid_time"`                     // 绑定时间（时间戳）
	Remark                    string          `gorm:"column:remark"`                       // 备注说明
	ShowId                    string          `gorm:"column:showId"`                       // 场次ID
	WechatNoticeNums          int             `gorm:"column:wechat_notice_nums"`           // 微信公众号通知成功数量
	TmallOrderNo              string          `gorm:"column:tmall_order_no"`               // 淘宝订单号
	TmallNickname             string          `gorm:"column:tmall_nickname"`               // 淘宝用户昵称（索引）
	PreCustomerServiceId      int64           `gorm:"column:pre_customer_service_id"`      // 售前客服ID
	AfterCustomerServiceId    int64           `gorm:"column:after_customer_service_id"`    // 售后客户ID
	RecordingOrderServiceName string          `gorm:"column:recording_order_service_name"` // 录单客服（索引）
	IsDirectGive              int8            `gorm:"column:is_direct_give"`               // 是否为一口价直接上报
	IsContractOrder           int8            `gorm:"column:isContractOrder"`              // 是否是承包订单
	BiddingBonusId            int64           `gorm:"column:bidding_bonus_id"`             // 票商结算ID
	BonusStatus               int8            `gorm:"column:bonus_status"`                 // 订单结算状态：1.已结算，2.未结算
	ShowTime                  int             `gorm:"column:show_time"`                    // 开场时间（时间戳）
	QuotationUseId            int64           `gorm:"column:quotation_use_id"`             // 自动报价ID（索引）
	CinemaSuggestRuleId       int64           `gorm:"column:cinema_suggest_rule_id"`       // 影院报价规则ID（索引）
	Dispute                   int8            `gorm:"column:dispute"`                      // 是否为确认纠纷订单：1是 2否
	ShopType                  int8            `gorm:"column:shop_type"`                    // 商家类型：1=淘宝
	CustomerId                string          `gorm:"column:customerId"`                   // 客户淘宝id（索引）
	MovieId                   int64           `gorm:"column:movieId"`                      // 电影ID（索引）
	HasPayment                int8            `gorm:"column:has_payment"`                  // 是否已付款：0未付款，1已付款，2已取消，3已退款
	CustomerUid               int64           `gorm:"column:customer_uid"`                 // 用户id（关联 customer_user.id）
	SuggestRulePriceId        int             `gorm:"column:suggest_rule_price_id"`        // 使用规则子分类价格id
	RefreshBiddingId          int             `gorm:"column:refresh_bidding_id"`           // 个人购票端刷新使用最终报价id
	EntryMethod               int8            `gorm:"column:entry_method"`                 // 录入方式：1.智能录单，2.手动录单
	TicketType                int8            `gorm:"column:ticket_type"`                  // 购票类型：1.竞价购票，2.特惠购票
	MerchantQuoteType         int8            `gorm:"column:merchant_quote_type"`          // 商家计价类型：1=成本价*比例，2=智能报价*比例
	LimitTime                 int             `gorm:"column:limit_time"`                   // 出票限制时间：0没有限制
	MerchantPrice             decimal.Decimal `gorm:"column:merchant_price"`               // 商家成本
	RetailerId                int64           `gorm:"column:retailer_id"`                  // 分销商ID
	SeqNo                     string          `gorm:"column:seqNo"`                        // 猫眼场次Id
}

func (Order) TableName() string {
	return "tb_order_jj"
}
