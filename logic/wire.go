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
	wire.Value(infra.DB),    // 提供DB
	wire.Value(infra.Redis), // 提供Redis
)

var svcSet = wire.NewSet(
	user.NewService,
	role.NewService,
	permission.NewService,
)

func initialize() *Service {
	wire.Build(wire.Struct(new(Service), "*"), providerSet, svcSet)
	return nil
}
