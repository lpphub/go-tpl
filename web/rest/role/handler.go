package role

import (
	"go-tpl/infra/logging"
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
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	data, err := logic.RoleSvc.List(c.Request.Context(), req)
	if err != nil {
		logging.Errorf(c, "Failed to get role list: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, data)
}

// Get 获取单个角色
func Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid role id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	role, err := logic.RoleSvc.Get(c.Request.Context(), uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to get role: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, role)
}

// Create 创建角色
func Create(c *gin.Context) {
	var req types.CreateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	role, err := logic.RoleSvc.Create(c.Request.Context(), req)
	if err != nil {
		logging.Errorf(c, "Failed to create role: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "Role created successfully: %d", role.ID)
	base.OKWithData(c, role)
}

// Update 更新角色
func Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid role id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.UpdateRoleReq
	if err = c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.RoleSvc.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		logging.Errorf(c, "Failed to update role: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "Role updated successfully: %d", id)
	base.OK(c)
}

// Delete 删除角色
func Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid role id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.RoleSvc.Delete(c.Request.Context(), uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to delete role: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "Role deleted successfully: %d", id)
	base.OK(c)
}

// UpdateStatus 更新角色状态
func UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid role id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.UpdateStatusReq
	if err = c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.RoleSvc.UpdateStatus(c.Request.Context(), uint(id), req.Status)
	if err != nil {
		logging.Errorf(c, "Failed to update role status: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "Role status updated successfully: %d, status: %d", id, req.Status)
	base.OK(c)
}

// GetRolePermissions 获取角色权限列表
func GetRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid role id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	permissionIds, err := logic.RoleSvc.GetRolePermissions(c.Request.Context(), uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to get role permissions: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, permissionIds)
}

// AssignPermissions 为角色分配权限
func AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid role id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.AssignRolePermissionsReq
	if err = c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.RoleSvc.AssignPermissions(c.Request.Context(), uint(id), req.PermissionIds)
	if err != nil {
		logging.Errorf(c, "Failed to assign permissions: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "Role permissions assigned successfully: %d, permissions: %v", id, req.PermissionIds)
	base.OK(c)
}

// GetRoleUsers 获取角色用户列表
func GetRoleUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid role id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	userIds, err := logic.RoleSvc.GetRoleUsers(c.Request.Context(), uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to get role users: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, userIds)
}
