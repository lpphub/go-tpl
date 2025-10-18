package user

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup) {
	r := router.Group("/user")
	{
		r.GET("/list", List)
		r.GET("/:id", GetByID)
	}
}
