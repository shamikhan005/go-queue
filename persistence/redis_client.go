package persistence

import (
	"context"
	"log"
	"time"

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

func (r *RedisClient) AddTaskWithState(queue string, taskID string, taskData string) error {
	pipe := r.Client.TxPipeline()

	pipe.RPush(ctx, queue, taskID)

	taskKey := "task:" + taskID
	pipe.HSet(ctx, taskKey, "data", taskData)
	pipe.HSet(ctx, taskKey, "status", "Pending")
	pipe.HSet(ctx, taskKey, "created_at", time.Now().Format(time.RFC3339))

	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisClient) UpdateTaskState(taskID string, state string) error {
	taskKey := "task:" + taskID
	return r.Client.HSet(ctx, taskKey, "status", string(state)).Err()
}

func (r *RedisClient) GetTaskData(taskID string) (map[string]string, error) {
	taskKey := "task:" + taskID
	return r.Client.HGetAll(ctx, taskKey).Result()
}