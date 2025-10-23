package rest

import (
	"go-tpl/infra/logger/logc"
	"go-tpl/web/base"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	logc.Info(c, "test1")

	logc.Info(c.Request.Context(), "test2")

	logc.Info(c, "test3")

	base.OKWithData(c, "ok")
}

func Register(c *gin.Context) {
	// todo  注册

	logc.Info(c, "register")
	base.OK(c)
}

func Login(c *gin.Context) {
	// todo  登录

	logc.Info(c, "login")

	base.OK(c)
}
