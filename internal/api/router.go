package api

import (
	"go-tpl/internal/api/handler"
	"go-tpl/internal/domain/repo"
	"go-tpl/internal/infra/global"
	"go-tpl/internal/service"

	"github.com/gin-gonic/gin"
)

type router struct {
	user *handler.UserHandler
}

func initAppContext() *router {
	// init repo
	var (
		userRepo = repo.NewUserRepo(global.DB)
	)

	// init service
	var (
		userSvc = service.NewUserService(userRepo, global.Redis)
	)

	return &router{
		user: handler.NewUserHandler(userSvc),
	}
}

func SetupRoute(app *gin.Engine) {

	route := initAppContext()

	u := app.Group("/user")
	{
		u.GET("/list", route.user.PageList)
	}
}
