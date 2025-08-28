package service

import (
	"go-tpl/internal/domain/entity"
	"go-tpl/internal/domain/repo"
	"go-tpl/pkg/pagination"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepo *repo.UserRepo
	redis    *redis.Client
}

func NewUserService(userRepo *repo.UserRepo, redis *redis.Client) *UserService {
	return &UserService{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (s *UserService) PageList(ctx *gin.Context, page pagination.Pagination) (*pagination.PageData[entity.User], error) {
	return s.userRepo.PageList(ctx, page)
}
