package entity

import (
	"go-tpl/pkg/ext"
)

type BaseModel struct {
	CreatedBy int64          `gorm:"column:created_by" json:"createdBy"`
	CreatedAt *ext.Timestamp `gorm:"column:created_at" json:"createdAt"`
	UpdatedBy int64          `gorm:"column:updated_by" json:"updatedBy"`
	UpdatedAt *ext.Timestamp `gorm:"column:updated_at" json:"updatedAt"`
	Deleted   int8           `gorm:"column:deleted;default:0" json:"deleted"`
}
