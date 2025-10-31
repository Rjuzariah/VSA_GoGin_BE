package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Docker service name
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}

	return rdb
}
