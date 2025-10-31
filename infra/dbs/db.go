package dbs

import (
	"context"
	"fmt"
	"go-tpl/infra/config"
	"go-tpl/infra/logging/logx"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlDB(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logx.NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	pool, _ := db.DB()
	pool.SetMaxIdleConns(2)
	pool.SetMaxOpenConns(10)
	pool.SetConnMaxIdleTime(5 * time.Minute)
	pool.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

func NewRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		MinIdleConns:    1,
		MaxActiveConns:  10,
		ConnMaxIdleTime: 3 * time.Minute,
		ConnMaxLifetime: 10 * time.Minute,
	})

	if err := client.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}

	client.AddHook(logx.NewRedisLogger())

	return client, nil
}
