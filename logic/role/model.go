package role

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Status      int8           `gorm:"default:1" json:"status"` // 1-正常, 0-禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Role) TableName() string {
	return "roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       uint `gorm:"column:role_id" json:"role_id"`
	PermissionID uint `gorm:"column:permission_id" json:"permission_id"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
