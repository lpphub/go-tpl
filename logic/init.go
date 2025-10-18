package logic

import (
	"go-tpl/infra"
	"go-tpl/logic/role"
	"go-tpl/logic/user"
)

var (
	UserSvc *user.UserService
	RoleSvc *role.RoleService
)

func Init() {
	UserSvc = user.NewUserService(infra.DB)
	RoleSvc = role.NewRoleService(infra.DB)
}
