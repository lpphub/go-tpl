package permission

import "github.com/gin-gonic/gin"

func Register(router *gin.RouterGroup) {
	r := router.Group("/permission")
	{
		r.POST("/list", List)                        // 获取权限列表
		r.GET("/:id", Get)                            // 获取单个权限
		r.POST("", Create)                            // 创建权限
		r.PUT("/:id", Update)                         // 更新权限
		r.DELETE("/:id", Delete)                      // 删除权限
		r.PUT("/:id/status", UpdateStatus)            // 更新权限状态
		r.GET("/modules", GetModules)                 // 获取所有模块
		r.GET("/:id/roles", GetPermissionRoles)       // 获取权限角色
	}
}
