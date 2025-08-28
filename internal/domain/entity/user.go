package entity

import "github.com/shopspring/decimal"

type User struct {
	UserId                  int64           `gorm:"column:user_id;primaryKey"`          // 主键ID
	Openid                  string          `gorm:"column:openid"`                      // 微信openid
	Unionid                 string          `gorm:"column:unionid"`                     // 应用unionid
	Nickname                string          `gorm:"column:nickname"`                    // 微信昵称
	MpId                    int64           `gorm:"column:mp_id"`                       // 公众号ID
	MpOpenid                string          `gorm:"column:mp_openid"`                   // 公众号openid
	AppOpenid               string          `gorm:"column:app_openid"`                  // 微信AppOpenid
	Level                   int             `gorm:"column:level"`                       // 会员等级
	Avatar                  string          `gorm:"column:avatar"`                      // 用户头像
	Gender                  int8            `gorm:"column:gender"`                      // 性别 2女1男
	Mobile                  string          `gorm:"column:mobile"`                      // 手机号码
	Email                   string          `gorm:"column:email"`                       // 邮箱
	RealName                string          `gorm:"column:real_name"`                   // 真实姓名
	IdCard                  string          `gorm:"column:idcard"`                      // 身份证号
	Password                string          `gorm:"column:password"`                    // 登陆密码
	Money                   decimal.Decimal `gorm:"column:money"`                       // 余额
	Points                  int             `gorm:"column:points"`                      // 积分
	BidPoints               int             `gorm:"column:bid_points"`                  // 抢单积分
	State                   int8            `gorm:"column:state"`                       // 0冻结1正常
	SuperiorUid             int64           `gorm:"column:superior_uid"`                // 推荐人ID
	SuperiorTime            int             `gorm:"column:superior_time"`               // 推荐绑定时间
	SuperiorTree            string          `gorm:"column:superior_tree"`               // 推荐关系树
	PCode                   string          `gorm:"column:pcode"`                       // 推荐码
	Qrcode                  string          `gorm:"column:qrcode"`                      // 二维码
	Status                  int8            `gorm:"column:status"`                      // 0 未认证 1 待审核 2 已通过 3拒绝
	LastLogin               int             `gorm:"column:last_login"`                  // 最后登陆时间
	Token                   string          `gorm:"column:token"`                       // 登录验证token
	RegTime                 int             `gorm:"column:reg_time"`                    // 注册时间
	TokenTime               int             `gorm:"column:token_time"`                  // 登录时间
	Type                    int8            `gorm:"column:type"`                        // 用户类型: 0用户,1竞价
	Exe                     int             `gorm:"column:exe"`                         // 经验值
	VipTime                 int             `gorm:"column:vip_time"`                    // vip到期时间
	From                    string          `gorm:"column:from"`                        // 注册渠道
	BiddingTime             string          `gorm:"column:bidding_time"`                // 自动竞价开始区间
	BiddingStatus           int8            `gorm:"column:bidding_status"`              // 自动竞价开关: 0关 ,1开
	IsOk                    int8            `gorm:"column:is_ok"`                       // 实名状态:0实名未通过,1实名通过
	IsFollow                int8            `gorm:"column:is_follow"`                   // 是否关注公众号
	IsPc                    int8            `gorm:"column:is_pc"`                       // 登录PC
	IsGz                    int8            `gorm:"column:is_gz"`                       // 是否关注 , 1关注 -1未关注
	Forbidden               int8            `gorm:"column:forbidden"`                   // 是否屏蔽订单通知
	ForbiddenTime           int             `gorm:"column:forbidden_time"`              // 屏蔽开始时间
	ForbiddenEndTime        string          `gorm:"column:forbidden_end_time"`          // 推送生效时间
	ContractSuspensionTime  int             `gorm:"column:contract_suspension_time"`    // 承包暂停时长
	IsContractAllocation    int8            `gorm:"column:is_contract_allocation"`      // 1.开启承包分配，2.关闭承包分配
	IsAutoForbidden         int8            `gorm:"column:is_auto_forbidden"`           // 是否自动关闭
	MaxBiddingTimes         int             `gorm:"column:max_bidding_times"`           // 最大同时竞价数量:-1代表取默认值
	ShowTutorials           int8            `gorm:"column:show_tutorials"`              // 接单教程，1.显示，2.不显示
	RedDotTips1             int8            `gorm:"column:red_dot_tips_1"`              // 个人中心红点提示，1提示，2.不提示
	RedDotTips2             int8            `gorm:"column:red_dot_tips_2"`              // 切换排版红点提示
	RedDotTips3             int8            `gorm:"column:red_dot_tips_3"`              // 哈哈助手红点提示
	BottomCover             int8            `gorm:"column:bottom_cover"`                // 1.兜底票商，2.NO
	IsQd                    int8            `gorm:"column:is_qd"`                       // 抢单权限，1.允许，2.否
	AuthAudit               int8            `gorm:"column:auth_audit"`                  // 认证审核，1.审核通过，2.审核未通过
	IsAutoBiddingLimitOrder int8            `gorm:"column:is_auto_bidding_limit_order"` // 是否自动报价限时单
}

func (User) TableName() string {
	return "tb_users"
}
