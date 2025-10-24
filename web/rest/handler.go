package rest

import (
	"go-tpl/infra/logging"
	"go-tpl/web/base"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	logging.Info(c, "test1")

	logging.Info(c.Request.Context(), "test2")

	logging.Info(c, "test3")

	base.OKWithData(c, "ok")
}

func Register(c *gin.Context) {
	// todo  注册

	logging.Info(c, "register")
	base.OK(c)
}

func Login(c *gin.Context) {
	// todo  登录

	logging.Info(c, "login")

	base.OK(c)
}
