package permission

import (
	"go-tpl/infra/logging"
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
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	data, err := logic.PermissionSvc.List(c.Request.Context(), req)
	if err != nil {
		logging.Errorf(c, "Failed to get permission list: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, data)
}

// Get 获取单个权限
func Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid permission id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	permission, err := logic.PermissionSvc.Get(c.Request.Context(), uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to get permission: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, permission)
}

// Create 创建权限
func Create(c *gin.Context) {
	var req types.CreatePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	permission, err := logic.PermissionSvc.Create(c.Request.Context(), req)
	if err != nil {
		logging.Errorf(c, "Failed to create permission: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Infof(c, "Permission created successfully: %d", permission.ID)
	base.OKWithData(c, permission)
}

// Update 更新权限
func Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid permission id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.UpdatePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.PermissionSvc.Update(c, uint(id), req)
	if err != nil {
		logging.Errorf(c, "Failed to update permission: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Infof(c, "Permission updated successfully: %d", id)
	base.OK(c)
}

// Delete 删除权限
func Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid permission id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.PermissionSvc.Delete(c, uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to delete permission: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Infof(c, "Permission deleted successfully: %d", id)
	base.OK(c)
}

// UpdateStatus 更新权限状态
func UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid permission id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.PermissionSvc.UpdateStatus(c, uint(id), req.Status)
	if err != nil {
		logging.Errorf(c, "Failed to update permission status: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Infof(c, "Permission status updated successfully: %d, status: %d", id, req.Status)
	base.OK(c)
}

// GetModules 获取所有模块列表
func GetModules(c *gin.Context) {
	modules, err := logic.PermissionSvc.GetModules(c)
	if err != nil {
		logging.Errorf(c, "Failed to get modules: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, modules)
}

// GetPermissionRoles 获取权限角色列表
func GetPermissionRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid permission id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	roleIds, err := logic.PermissionSvc.GetPermissionRoles(c, uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to get permission roles: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, roleIds)
}
