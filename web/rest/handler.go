package rest

import (
	"go-tpl/web/base"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	base.OKWithData(c, "ok")
}

func Register(c *gin.Context) {
	// todo  注册
	base.OK(c)
}

func Login(c *gin.Context) {
	// todo  登录
	base.OK(c)
}
