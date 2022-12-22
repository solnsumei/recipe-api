package services

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v9"
)

func InitRedis(ctx context.Context) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)

	return redisClient
}
