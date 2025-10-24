package user

import (
	"go-tpl/infra/logging"
	"go-tpl/logic"
	"go-tpl/logic/shared"
	"go-tpl/web/base"
	"go-tpl/web/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List 获取用户列表
func List(c *gin.Context) {
	var req types.UserQueryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	data, err := logic.UserSvc.List(c.Request.Context(), req)
	if err != nil {
		logging.Errorf(c, "Failed to get user list: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, data)
}

// Get 获取单个用户
func Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid user id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	user, err := logic.UserSvc.Get(c.Request.Context(), uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to get user: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, user)
}

// Create 创建用户
func Create(c *gin.Context) {
	var req types.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	user, err := logic.UserSvc.Create(c.Request.Context(), req)
	if err != nil {
		logging.Errorf(c, "Failed to create user: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "User created successfully: %d", user.ID)
	base.OKWithData(c, user)
}

// Update 更新用户
func Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid user id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.UpdateUserReq
	if err = c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.UserSvc.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		logging.Errorf(c, "Failed to update user: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "User updated successfully: %d", id)
	base.OK(c)
}

// Delete 删除用户
func Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid user id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.UserSvc.Delete(c.Request.Context(), uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to delete user: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "User deleted successfully: %d", id)
	base.OK(c)
}

// UpdateStatus 更新用户状态
func UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid user id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.UpdateStatusReq
	if err = c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.UserSvc.UpdateStatus(c.Request.Context(), uint(id), req.Status)
	if err != nil {
		logging.Errorf(c, "Failed to update user status: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "User status updated successfully: %d, status: %d", id, req.Status)
	base.OK(c)
}

// GetUserRoles 获取用户角色列表
func GetUserRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid user id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	roleIds, err := logic.UserSvc.GetUserRoles(c.Request.Context(), uint(id))
	if err != nil {
		logging.Errorf(c, "Failed to get user roles: %v", err)
		base.FailWithError(c, err)
		return
	}

	base.OKWithData(c, roleIds)
}

// AssignRoles 为用户分配角色
func AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logging.Errorf(c, "Invalid user id: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	var req types.AssignUserRolesReq
	if err = c.ShouldBindJSON(&req); err != nil {
		logging.Errorf(c, "Invalid request: %v", err)
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	err = logic.UserSvc.AssignRoles(c.Request.Context(), uint(id), req.RoleIds)
	if err != nil {
		logging.Errorf(c, "Failed to assign roles: %v", err)
		base.FailWithError(c, err)
		return
	}

	logging.Errorf(c, "User roles assigned successfully: %d, roles: %v", id, req.RoleIds)
	base.OK(c)
}
