package user

import (
	"context"
	"go-tpl/logic/shared"
	"go-tpl/web/types"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewService(db *gorm.DB, redis *redis.Client) *Service {
	return &Service{
		db:    db,
		redis: redis,
	}
}

func (s *Service) List(ctx context.Context, req types.UserQueryReq) (*shared.PageData[User], error) {
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
		if err := _db.Scopes(shared.Paginate(req.Pagination)).Find(&list).Error; err != nil {
			return nil, err
		}
	}
	return shared.Wrapper[User](total, list), nil
}
