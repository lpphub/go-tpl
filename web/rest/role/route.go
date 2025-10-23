package role

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup) {
	r := router.Group("/role")
	{
		r.POST("/list", List)                           // 获取角色列表
		r.GET("/:id", Get)                               // 获取单个角色
		r.POST("", Create)                               // 创建角色
		r.PUT("/:id", Update)                            // 更新角色
		r.DELETE("/:id", Delete)                         // 删除角色
		r.PUT("/:id/status", UpdateStatus)               // 更新角色状态
		r.GET("/:id/permissions", GetRolePermissions)    // 获取角色权限
		r.PUT("/:id/permissions", AssignPermissions)     // 分配角色权限
		r.GET("/:id/users", GetRoleUsers)                // 获取角色用户
	}
}
