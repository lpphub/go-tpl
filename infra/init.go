package infra

import (
	"fmt"
	"go-tpl/infra/config"
	"go-tpl/infra/dbs"
	"go-tpl/infra/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Cfg    *config.Config
	Logger *zap.Logger
	DB     *gorm.DB
	Redis  *redis.Client
)

func Init() {
	var err error
	// 1.加载配置
	Cfg, err = config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 2.初始化日志
	Logger = logger.NewLogger("")

	// 3.初始化数据库和Redis
	DB, err = dbs.NewMysqlDB(Cfg.Database)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	Redis, err = dbs.NewRedis(Cfg.Redis)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to redis: %v", err))
	}

}
