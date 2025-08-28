package entity

import "github.com/shopspring/decimal"

type BiddingBill struct {
	Id          int64           `gorm:"column:id;primaryKey"` // 主键ID（自增）
	OrderId     int64           `gorm:"column:order_id"`      // 关联订单ID
	Uid         int64           `gorm:"column:uid"`           // 出价人ID（关联用户表）
	Price       decimal.Decimal `gorm:"column:price"`         // 出价金额
	IsSeat      int8            `gorm:"column:is_seat"`       // 是否允许调座：0不允许，1允许
	LoverSeat   int8            `gorm:"column:lover_seat"`    // 是否情侣座：0否，1是
	ToSeat      int8            `gorm:"column:to_seat"`       // 是否已调座：0未调座，1已调座
	OrderSeat   string          `gorm:"column:order_seat"`    // 订单原始座位信息
	BindSeat    string          `gorm:"column:bind_seat"`     // 调位后的座位信息
	Img         string          `gorm:"column:img"`           // 凭证图片（文本存储，如URL或Base64）
	TicketCode  string          `gorm:"column:ticket_code"`   // 取票码（文本，支持多码存储）
	VerCode     string          `gorm:"column:ver_code"`      // 验证码
	OcrCode     string          `gorm:"column:ocr_code"`      // OCR识别后的取票码
	Remark      string          `gorm:"column:remark"`        // 备注信息
	Text        string          `gorm:"column:text"`          // 弃单相关说明
	Status      int8            `gorm:"column:status"`        // 竞价状态：0竞价中,1成功,2失败,3待上报,4待结算,5已完成,6已关闭,7纠纷中
	IsSent      int8            `gorm:"column:is_sent"`       // 剩余5分钟通知标记：0未通知，1已通知
	Nopay       int8            `gorm:"column:nopay"`         // 取消标记：0未取消，1已取消
	IsOver      int8            `gorm:"column:is_over"`       // 超时订单标记：0未超时，1已超时
	IsSee       int8            `gorm:"column:is_see"`        // 查看标记（冗余字段，与see_time配合）
	IsAuto      int8            `gorm:"column:is_auto"`       // 自动竞价标记：0手动竞价，1自动竞价
	IsCall      int8            `gorm:"column:is_call"`       // 电话沟通标记：0未电话，1已电话
	AutoBill    int8            `gorm:"column:auto_bill"`     // 自动分配标记：0手动分配，1自动分配
	AutoId      int64           `gorm:"column:auto_id"`       // 关联自动竞价规则ID
	IsQd        int8            `gorm:"column:is_qd"`         // 是否抢单：0非抢单，1抢单
	Edit        int8            `gorm:"column:edit"`          // 修改短信标记：0未修改，1已修改
	FromId      string          `gorm:"column:fromId"`        // 渠道标识
	SetTime     int             `gorm:"column:set_time"`      // 结算时间（时间戳）
	BindTime    int             `gorm:"column:bindd_time"`    // 绑定订单时间（时间戳）
	CreateTime  int             `gorm:"column:create_time"`   // 竞价记录生成时间（时间戳）
	ErTime      int             `gorm:"column:er_time"`       // 竞价失败时间（时间戳）
	UpdateTime  int             `gorm:"column:update_time"`   // 竞价状态更新时间（时间戳）
	SeeTime     int             `gorm:"column:see_time"`      // 查看时间（时间戳）
	AddCount    int             `gorm:"column:add_count"`     // 竞价时间延长次数
	AddTime     int             `gorm:"column:add_time"`      // 延长的时间（单位通常为秒）
	FeeCount    int             `gorm:"column:fee_count"`     // 费用计算标记（业务自定义）
	FeeAdd      int             `gorm:"column:fee_add"`       // 后台手动延长时间标记：0未延长，1已延长
	FirstTime   int             `gorm:"column:first_time"`    // 第一次查看时间（时间戳）
	DeleteTime  int             `gorm:"column:delete_time"`   // 记录删除时间（时间戳，软删除用）
	DelCount    int             `gorm:"column:del_count"`     // 删除次数（软删除重试标记）
	UCount      int             `gorm:"column:u_count"`       // 短信发送次数
	IsUa        int8            `gorm:"column:is_ua"`         // 承包订单标记：0非承包，1承包
	IsAdmin     int8            `gorm:"column:is_admin"`      // 后台手动失败标记：0非后台操作，1后台操作
	UploadTime  int             `gorm:"column:upload_time"`   // 竞价成功提交时间（时间戳）
	IsConfirm   int8            `gorm:"column:is_confirm"`    // 订单确认状态（业务自定义，如1已确认）
	ConfirmTime int             `gorm:"column:confirm_time"`  // 订单确认时间（时间戳）
	CinemaId    int64           `gorm:"column:cinema_id"`     // 关联影院ID
	ImgBak      string          `gorm:"column:img_bak"`       // 凭证图片备份（冗余字段）
}

func (BiddingBill) TableName() string {
	return "tb_bidding_bill"
}
