package api

import (
	"go-tpl/server/infra"

	"github.com/gin-gonic/gin"
)

func SetupRoute(engine *gin.Engine) {

	app := infra.InitAppContext()

	u := engine.Group("/user")
	{
		u.GET("/list", app.User.PageList)
	}
}
