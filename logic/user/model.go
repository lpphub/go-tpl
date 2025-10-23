package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"size:100;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Status    int8           `gorm:"default:1" json:"status"` // 1-正常, 0-禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID uint `gorm:"column:user_id" json:"user_id"`
	RoleID uint `gorm:"column:role_id" json:"role_id"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
