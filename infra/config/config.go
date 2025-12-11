package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret            string
	ExpireTime        int64 // access_token 过期时间（秒）
	RefreshExpireTime int64 // refresh_token 过期时间（秒）
}

type ServerConfig struct {
	Port int
	Mode string // debug, release, test
}

type Config struct {
	Database DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Server   ServerConfig
}

func Load() (*Config, error) {
	return load[Config]("./config", "conf", "yml")
}

func load[T any](configPath, configName, configType string) (*T, error) {
	v := viper.New()

	// 1. 配置文件基础设置
	v.AddConfigPath(configPath) // 设置配置文件目录
	v.SetConfigName(configName) // 设置配置文件名称（无后缀）
	v.SetConfigType(configType) // 设置配置文件类型（显式指定，兼容无后缀文件）

	// 2. 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		// 忽略配置文件不存在的错误（允许仅通过环境变量配置）
		var cffErr viper.ConfigFileNotFoundError
		if !errors.As(err, &cffErr) {
			return nil, fmt.Errorf("read config file failed: %w", err)
		}
	}

	// 3. 环境变量配置（优先级高于配置文件）
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	var cfg T
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config to struct failed: %w", err)
	}

	return &cfg, nil
}
