package types

import "go-tpl/logic/shared"

// 用户相关请求类型
type UserQueryReq struct {
	shared.Pagination
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   *int8  `json:"status"`
}

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   *int8  `json:"status"`
}

// 角色相关请求类型
type RoleQueryReq struct {
	shared.Pagination
	Name   string `json:"name"`
	Status *int8  `json:"status"`
}

type CreateRoleReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleReq struct {
	Name        string `json:"name"`
	Description *string `json:"description"`
	Status      *int8  `json:"status"`
}

type AssignRolePermissionsReq struct {
	PermissionIds []uint `json:"permission_ids" binding:"required"`
}

// 权限相关请求类型
type PermissionQueryReq struct {
	shared.Pagination
	Code   string `json:"code"`
	Name   string `json:"name"`
	Module string `json:"module"`
	Status *int8  `json:"status"`
}

type CreatePermissionReq struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Module      string `json:"module" binding:"required"`
}

type UpdatePermissionReq struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description *string `json:"description"`
	Module      *string `json:"module"`
	Status      *int8   `json:"status"`
}

// 用户角色相关请求类型
type AssignUserRolesReq struct {
	RoleIds []uint `json:"role_ids" binding:"required"`
}

// 状态更新请求类型
type UpdateStatusReq struct {
	Status int8 `json:"status" binding:"required"`
}

// 登录注册请求类型
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
