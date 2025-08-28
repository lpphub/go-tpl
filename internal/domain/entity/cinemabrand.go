package entity

import (
	"go-tpl/pkg/ext"
	"go-tpl/pkg/util"
)

type CinemaBrand struct {
	Id         int64          `gorm:"column:id;primaryKey"` // 主键ID
	Name       string         `gorm:"column:name"`          // 影院品牌名称
	Status     int            `gorm:"column:status"`        // 品牌状态：1启用
	CinemaIds  string         `gorm:"column:cinemaids"`     // 该品牌下的影院ID集合
	CreateDate *ext.Timestamp `gorm:"column:create_date"`   // 品牌创建时间
}

func (CinemaBrand) TableName() string {
	return "tb_hh_cinema_brand"
}

func (e CinemaBrand) SplitCinemaIds() []int64 {
	if e.CinemaIds == "" {
		return nil
	}
	return util.SplitToInt64Slice(e.CinemaIds, ",")
}
