package permission

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Module      string         `gorm:"size:50;index" json:"module"`
	Status      int8           `gorm:"default:1" json:"status"` // 1-正常, 0-禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
	return "permissions"
}
