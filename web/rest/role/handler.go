package role

import (
	"go-tpl/infra/logger"
	"go-tpl/logic"
	"go-tpl/logic/shared"
	"go-tpl/web/base"
	"go-tpl/web/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List 获取角色列表
func List(c *gin.Context) {
	var req types.RoleQueryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	data, err := logic.Svc.Role.List(c, req)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, data)
}

// Get 获取单个角色
func Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	role, err := logic.Svc.Role.Get(c, uint(id))
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, role)
}

// Create 创建角色
func Create(c *gin.Context) {
	var req types.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	role, err := logic.Svc.Role.Create(c, req)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, role)
}

// Update 更新角色
func Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.UpdateRoleReq
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.Svc.Role.Update(c, uint(id), req)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OK(c)
}

// Delete 删除角色
func Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.Svc.Role.Delete(c, uint(id))
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OK(c)
}

// UpdateStatus 更新角色状态
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

	err = logic.Svc.Role.UpdateStatus(c, uint(id), req.Status)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OK(c)
}

// GetRolePermissions 获取角色权限列表
func GetRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	permissionIds, err := logic.Svc.Role.GetRolePermissions(c, uint(id))
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, permissionIds)
}

// AssignPermissions 为角色分配权限
func AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.AssignRolePermissionsReq
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.Svc.Role.AssignPermissions(c, uint(id), req.PermissionIds)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OK(c)
}

// GetRoleUsers 获取角色用户列表
func GetRoleUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	userIds, err := logic.Svc.Role.GetRoleUsers(c, uint(id))
	if err != nil {
		logger.Errw(c, err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, userIds)
}
