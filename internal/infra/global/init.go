package global

import (
	"go-tpl/internal/infra/conf"

	"github.com/lpphub/golib/env"
	"github.com/lpphub/golib/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Conf  conf.RConfig
	DB    *gorm.DB
	Redis *redis.Client
)

func preInit() {
	env.SetAppName("flow-flow")

	Conf = conf.LoadConfig()

	logger.Setup()
}

func InitResource() {
	preInit()

	initDb()
	initRedis()
}

func Clear() {
}
