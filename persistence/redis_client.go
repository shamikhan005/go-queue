package persistence

import (
	"context"
	"log"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	return &RedisClient{Client: client}
}

func (r *RedisClient) AddTask(queue string, taskData string) error {
	return r.Client.RPush(ctx, queue, taskData).Err()
}

func (r *RedisClient) GetTask(queue string) (string, error) {
	return r.Client.LPop(ctx, queue).Result()
}