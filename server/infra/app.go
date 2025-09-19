package infra

import (
	"go-tpl/server/api/handler"
	"go-tpl/server/domain/repo"
	"go-tpl/server/infra/global"
	"go-tpl/server/service"
)

type AppContext struct {
	User *handler.UserHandler
}

// InitAppContext 初始化应用上下文
func InitAppContext() *AppContext {
	// init repo
	repoApp := initRepoContainer()
	// init service
	svcApp := initServiceContainer(repoApp)

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

func initRepoContainer() *repoContainer {
	return &repoContainer{
		UserRepo: repo.NewUserRepo(global.DB),
	}
}

func initServiceContainer(repo *repoContainer) *serviceContainer {
	return &serviceContainer{
		UserSvc: service.NewUserService(repo.UserRepo, global.Redis),
	}
}
