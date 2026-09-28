package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}
