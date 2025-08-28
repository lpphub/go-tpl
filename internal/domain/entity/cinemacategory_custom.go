package entity

import (
	"go-tpl/pkg/ext"
	"go-tpl/pkg/util"
)

type UserCustomCinemaCategory struct {
	Id         int64          `gorm:"column:id;primaryKey"` // 主键ID
	Uid        int64          `gorm:"column:uid"`           // 用户id
	Name       string         `gorm:"column:name"`          // 规则名称
	CinemaId   string         `gorm:"column:cinemaId"`      // 影院ids（多ID用文本存储）
	DelTime    int            `gorm:"column:delTime"`       // 删除时间（时间戳）
	Status     int            `gorm:"column:status"`        // 状态：1=开，2=关
	CreateTime *ext.Timestamp `gorm:"column:createTime"`    // 创建时间
	UpdateTime *ext.Timestamp `gorm:"column:updateTime"`    // 更新时间
}

func (UserCustomCinemaCategory) TableName() string {
	return "tb_hh_user_custom_cinema_category"
}

func (c UserCustomCinemaCategory) SplitCinemaIds() []int64 {
	return util.SplitToInt64Slice(c.CinemaId, ",")
}
