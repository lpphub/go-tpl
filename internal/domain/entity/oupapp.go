package entity

import "go-tpl/pkg/ext"

type OupAppSecret struct {
	Id          int64          `gorm:"column:id;primaryKey"` // 主键ID（数据库为 int unsigned，对应Go的uint64）
	UserId      int64          `gorm:"column:userId"`        // 票商ID（关联用户表，唯一标识票商，唯一索引）
	AppId       string         `gorm:"column:appId"`         // 接口应用ID（票商调用接口的唯一标识）
	AppSecret   string         `gorm:"column:appSecret"`     // 接口应用密钥（票商调用接口的鉴权密钥）
	NoticeUrl   string         `gorm:"column:noticeUrl"`     // 通知URL（接口回调通知地址，如订单状态变更通知）
	UriAuthList string         `gorm:"column:uriAuthList"`   // 接口权限列表（文本存储，如逗号分隔的接口URI，控制票商可调用的接口范围）
	Status      int            `gorm:"column:status"`        // 权限状态：1开启（允许调用），其他值可表示禁用（参考数据库默认值1）
	CreateTime  *ext.Timestamp `gorm:"column:createTime"`    // 记录创建时间
	UpdateTime  *ext.Timestamp `gorm:"column:updateTime"`    // 记录更新时间
}

func (OupAppSecret) TableName() string {
	return "tb_hh_oup_app_secret"
}
