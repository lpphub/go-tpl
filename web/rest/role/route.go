package role

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup) {
	r := router.Group("/role")
	{
		r.GET("/list", List)
		r.GET("/:id", GetByID)
	}
}
