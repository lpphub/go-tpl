package flow

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Storage interface {
	GetCount(ctx context.Context, setKey string) (int64, error)
	GetMembers(ctx context.Context, setKey string) ([]string, error)
}

type RedisStorage struct {
	Redis *redis.Client
}

func (r RedisStorage) GetCount(ctx context.Context, setKey string) (int64, error) {
	//TODO implement me
	panic("implement me")

}

func (r RedisStorage) GetMembers(ctx context.Context, setKey string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}
