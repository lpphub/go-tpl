package web

import (
	"go-tpl/infra/logging/logx"
	"go-tpl/web/rest"
	"go-tpl/web/rest/permission"
	"go-tpl/web/rest/role"
	"go-tpl/web/rest/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(logx.GinLogMiddleware())

	api := r.Group("/api")

	// 公共路由
	api.GET("/test", rest.Test)
	api.POST("/register", rest.Register)
	api.POST("/login", rest.Login)

	auth := api.Group("")
	{
		// 注册接口处理器
		user.Register(auth)
		role.Register(auth)
		permission.Register(auth)
	}

	return r
}
