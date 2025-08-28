package entity

type User struct {
	ID       int64  `gorm:"column:id"`       // 主键ID
	Nickname string `gorm:"column:nickname"` // 昵称
	BaseModel
}

func (User) TableName() string {
	return "tb_user"
}
