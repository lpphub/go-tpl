package user

import (
	"context"
	"go-tpl/logic/base"
	"go-tpl/web/types"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) List(ctx context.Context, req types.UserQueryReq) (*base.PageData[User], error) {
	var (
		total int64
		list  []User
	)

	_db := s.db.WithContext(ctx).Model(&User{})
	if req.Username != "" {
		_db.Where("username like ?", "%"+req.Username+"%")
	}
	if err := _db.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		if err := _db.Scopes(base.Paginate(req.Pagination)).Find(&list).Error; err != nil {
			return nil, err
		}
	}
	return base.Wrapper[User](total, list), nil
}
