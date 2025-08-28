package global

import (
	"context"
	"go-tpl/pkg/ext/logext"
	"time"

	"github.com/redis/go-redis/v9"
)

func initRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:            Conf.Redis.Addr,
		Password:        Conf.Redis.Password,
		DB:              Conf.Redis.DB,
		MinIdleConns:    8,
		MaxActiveConns:  20,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 30 * time.Second,
	})
	Redis.AddHook(logext.RedisLogHook{})

	if _, err := Redis.Ping(context.Background()).Result(); err != nil {
		panic("init redis error: " + err.Error())
	}
}
