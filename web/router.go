package web

import (
	"go-tpl/web/rest"
	"go-tpl/web/rest/role"
	"go-tpl/web/rest/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	// 公共路由
	api.GET("/test", rest.Test)

	api.POST("/register", rest.Register)
	api.POST("/login", rest.Login)

	//
	//// 注册接口路由
	user.Register(api)
	role.Register(api)
	//permission.Register(api)
	return r
}
