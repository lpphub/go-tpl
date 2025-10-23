package permission

import (
	"go-tpl/infra/logger/logc"
	"go-tpl/logic"
	"go-tpl/web/base"
	"go-tpl/web/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List 获取权限列表
func List(c *gin.Context) {
	var req types.PermissionQueryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logc.Errorf(c.Request.Context(), "Invalid request: %v", err)
		base.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	data, err := logic.PermissionSvc.List(c.Request.Context(), req)
	if err != nil {
		logc.Errorf(c.Request.Context(), "Failed to get permission list: %v", err)
		base.FailWithErr(c, err)
		return
	}

	base.OKWithData(c, data)
}

// Get 获取单个权限
func Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logc.Errorf(c.Request.Context(), "Invalid permission id: %v", err)
		base.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	permission, err := logic.PermissionSvc.Get(c.Request.Context(), uint(id))
	if err != nil {
		logc.Errorf(c.Request.Context(), "Failed to get permission: %v", err)
		base.FailWithErr(c, err)
		return
	}

	base.OKWithData(c, permission)
}

// Create 创建权限
func Create(c *gin.Context) {
	var req types.CreatePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logc.Errorf(c.Request.Context(), "Invalid request: %v", err)
		base.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	permission, err := logic.PermissionSvc.Create(c.Request.Context(), req)
	if err != nil {
		logc.Errorf(c.Request.Context(), "Failed to create permission: %v", err)
		base.FailWithErr(c, err)
		return
	}

	logc.Infof(c.Request.Context(), "Permission created successfully: %d", permission.ID)
	base.OKWithData(c, permission)
}

// Update 更新权限
func Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logc.Errorf(c.Request.Context(), "Invalid permission id: %v", err)
		base.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	var req types.UpdatePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logc.Errorf(c.Request.Context(), "Invalid request: %v", err)
		base.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	err = logic.PermissionSvc.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		logc.Errorf(c.Request.Context(), "Failed to update permission: %v", err)
		base.FailWithErr(c, err)
		return
	}

	logc.Infof(c.Request.Context(), "Permission updated successfully: %d", id)
	base.OK(c)
}

// Delete 删除权限
func Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logc.Errorf(c.Request.Context(), "Invalid permission id: %v", err)
		base.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	err = logic.PermissionSvc.Delete(c.Request.Context(), uint(id))
	if err != nil {
		logc.Errorf(c.Request.Context(), "Failed to delete permission: %v", err)
		base.FailWithErr(c, err)
		return
	}

	logc.Infof(c.Request.Context(), "Permission deleted successfully: %d", id)
	base.OK(c)
}

// UpdateStatus 更新权限状态
func UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logc.Errorf(c.Request.Context(), "Invalid permission id: %v", err)
		base.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	var req types.UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logc.Errorf(c.Request.Context(), "Invalid request: %v", err)
		base.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	err = logic.PermissionSvc.UpdateStatus(c.Request.Context(), uint(id), req.Status)
	if err != nil {
		logc.Errorf(c.Request.Context(), "Failed to update permission status: %v", err)
		base.FailWithErr(c, err)
		return
	}

	logc.Infof(c.Request.Context(), "Permission status updated successfully: %d, status: %d", id, req.Status)
	base.OK(c)
}

// GetModules 获取所有模块列表
func GetModules(c *gin.Context) {
	modules, err := logic.PermissionSvc.GetModules(c.Request.Context())
	if err != nil {
		logc.Errorf(c.Request.Context(), "Failed to get modules: %v", err)
		base.FailWithErr(c, err)
		return
	}

	base.OKWithData(c, modules)
}

// GetPermissionRoles 获取权限角色列表
func GetPermissionRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logc.Errorf(c.Request.Context(), "Invalid permission id: %v", err)
		base.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	roleIds, err := logic.PermissionSvc.GetPermissionRoles(c.Request.Context(), uint(id))
	if err != nil {
		logc.Errorf(c.Request.Context(), "Failed to get permission roles: %v", err)
		base.FailWithErr(c, err)
		return
	}

	base.OKWithData(c, roleIds)
}
