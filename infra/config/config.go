package config

import (
	"os"
	"strconv"

	"github.com/goccy/go-yaml"
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
	Secret     string
	ExpireTime int64 // 过期时间（秒）
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
	// 首先尝试从配置文件加载
	cfgFile, err := loadFromFile("config/conf.yml")
	if err != nil {
		// 如果配置文件不存在或解析失败，使用默认配置
		cfgFile = &Config{}
	}

	// 环境变量优先，覆盖配置文件中的值
	return &Config{
		Database: DBConfig{
			Host:     getEnv("DB_HOST", cfgFile.Database.Host),
			Port:     getEnvAsInt("DB_PORT", cfgFile.Database.Port),
			User:     getEnv("DB_USER", cfgFile.Database.User),
			Password: getEnv("DB_PASSWORD", cfgFile.Database.Password),
			Dbname:   getEnv("DB_NAME", cfgFile.Database.Dbname),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", cfgFile.Redis.Host),
			Port:     getEnvAsInt("REDIS_PORT", cfgFile.Redis.Port),
			Password: getEnv("REDIS_PASSWORD", cfgFile.Redis.Password),
			DB:       getEnvAsInt("REDIS_DB", cfgFile.Redis.DB),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", cfgFile.JWT.Secret),
			ExpireTime: getEnvAsInt64("JWT_EXPIRE_TIME", cfgFile.JWT.ExpireTime), // 默认24小时
		},
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", cfgFile.Server.Port),
			Mode: getEnv("SERVER_MODE", cfgFile.Server.Mode),
		},
	}, nil
}

// loadFromFile 从YAML配置文件加载配置
func loadFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		value, err := strconv.Atoi(valueStr)
		if err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if valueStr := os.Getenv(key); valueStr != "" {
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err == nil {
			return value
		}
	}
	return defaultValue
}
