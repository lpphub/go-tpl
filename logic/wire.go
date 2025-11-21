//go:build wireinject
// +build wireinject

//go:generate wire

package logic

import (
	"go-tpl/infra"
	"go-tpl/logic/permission"
	"go-tpl/logic/role"
	"go-tpl/logic/user"

	"github.com/google/wire"
)

type Service struct {
	User       *user.Service
	Role       *role.Service
	Permission *permission.Service
}

var providerSet = wire.NewSet(
	infra.ProvideDB,  // 提供DB
	infra.ProvideRDB, // 提供Redis
)

var svcSet = wire.NewSet(
	user.NewService,
	role.NewService,
	permission.NewService,
)

func initialize() *Service {
	wire.Build(providerSet, svcSet, wire.Struct(new(Service), "*"))
	return nil
}
