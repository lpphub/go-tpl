package role

import (
	"context"
	"go-tpl/logic/shared"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) List(ctx context.Context, page shared.Pagination) (*shared.PageData[Role], error) {
	var (
		total int64
		list  []Role
	)
	_db := s.db.WithContext(ctx).Model(&Role{})
	if err := _db.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		if err := _db.Scopes(shared.Paginate(page)).Find(&list).Error; err != nil {
			return nil, err
		}
	}
	return shared.Wrapper[Role](total, list), nil
}
