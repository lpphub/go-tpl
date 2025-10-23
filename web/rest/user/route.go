package user

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup) {
	r := router.Group("/user")
	{
		r.POST("/list", List)        // 获取用户列表
		r.GET("/:id", Get)            // 获取单个用户
		r.POST("", Create)            // 创建用户
		r.PUT("/:id", Update)         // 更新用户
		r.DELETE("/:id", Delete)      // 删除用户
		r.PUT("/:id/status", UpdateStatus) // 更新用户状态
		r.GET("/:id/roles", GetUserRoles)   // 获取用户角色
		r.PUT("/:id/roles", AssignRoles)    // 分配用户角色
	}
}
