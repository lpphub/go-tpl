package permission

import (
	"go-tpl/infra/logger"
	"go-tpl/logic"
	"go-tpl/logic/shared"
	"go-tpl/web/base"
	"go-tpl/web/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List 获取权限列表
func List(c *gin.Context) {
	var req types.PermissionQueryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	data, err := logic.Svc.Permission.List(c, req)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, data)
}

// Get 获取单个权限
func Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	permission, err := logic.Svc.Permission.Get(c, uint(id))
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, permission)
}

// Create 创建权限
func Create(c *gin.Context) {
	var req types.CreatePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	permission, err := logic.Svc.Permission.Create(c, req)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, permission)
}

// Update 更新权限
func Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.UpdatePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.Svc.Permission.Update(c, uint(id), req)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OK(c)
}

// Delete 删除权限
func Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.Svc.Permission.Delete(c, uint(id))
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OK(c)
}

// UpdateStatus 更新权限状态
func UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.UpdateStatusReq
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.Svc.Permission.UpdateStatus(c, uint(id), req.Status)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OK(c)
}

// GetModules 获取所有模块列表
func GetModules(c *gin.Context) {
	modules, err := logic.Svc.Permission.GetModules(c)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, modules)
}

// GetPermissionRoles 获取权限角色列表
func GetPermissionRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	roleIds, err := logic.Svc.Permission.GetPermissionRoles(c, uint(id))
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, roleIds)
}
