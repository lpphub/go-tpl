package role

import (
	"context"
	"go-tpl/logic/base"

	"gorm.io/gorm"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		db: db,
	}
}

func (s *RoleService) List(ctx context.Context, page base.Pagination) (*base.PageData[Role], error) {
	var (
		total int64
		list  []Role
	)
	_db := s.db.WithContext(ctx).Model(&Role{})
	if err := _db.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		if err := _db.Scopes(base.Paginate(page)).Find(&list).Error; err != nil {
			return nil, err
		}
	}
	return base.Wrapper[Role](total, list), nil
}
