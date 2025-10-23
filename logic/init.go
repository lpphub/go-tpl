package logic

import (
	"go-tpl/infra"
	"go-tpl/logic/permission"
	"go-tpl/logic/role"
	"go-tpl/logic/user"
)

var (
	UserSvc       *user.Service
	RoleSvc       *role.Service
	PermissionSvc *permission.Service
)

func Init() {
	UserSvc = user.NewService(infra.DB, infra.Redis)
	RoleSvc = role.NewService(infra.DB)
	PermissionSvc = permission.NewService(infra.DB)
}
