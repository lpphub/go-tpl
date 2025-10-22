package permission

import "github.com/gin-gonic/gin"

func Register(router *gin.RouterGroup) {
	r := router.Group("/permission")
	{
		r.GET("/list", List)
		r.GET("/:id", GetByID)
	}
}
