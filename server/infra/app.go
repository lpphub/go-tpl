package infra

import (
	"go-tpl/server/api/handler"
	"go-tpl/server/domain/repo"
	"go-tpl/server/infra/global"
	"go-tpl/server/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppContext struct {
	User *handler.UserHandler
}

// InitAppContext 初始化应用上下文
func InitAppContext() *AppContext {
	// init repo
	repoApp := initRepoContainer(global.DB)
	// init service
	svcApp := initServiceContainer(repoApp, global.Redis)

	return &AppContext{
		User: handler.NewUserHandler(svcApp.UserSvc),
	}
}

type serviceContainer struct {
	UserSvc *service.UserService
}

type repoContainer struct {
	UserRepo *repo.UserRepo
}

func initRepoContainer(db *gorm.DB) *repoContainer {
	return &repoContainer{
		UserRepo: repo.NewUserRepo(db),
	}
}

func initServiceContainer(repo *repoContainer, redis *redis.Client) *serviceContainer {
	return &serviceContainer{
		UserSvc: service.NewUserService(repo.UserRepo, redis),
	}
}
