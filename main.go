package main

import (
	"context"
	"encoding/json"
	"first/go-queue/persistence"
	"first/go-queue/queue"
	"first/go-queue/task"
	"flag"
	"fmt"
	"log"
)

func main() {
	taskName := flag.String("task", "default-task", "name of the task to enqueue")
	flag.Parse()

	redisClient := persistence.NewRedisClient()
	ctx := context.Background()

	newTask := task.NewTask(*taskName)

	taskData, err := json.Marshal(newTask)
	if err != nil {
		log.Fatalf("failed to serialize task: %v", err)
	}

	queueName := "task-queue"
	err = redisClient.AddTaskWithState(queueName, newTask.ID, string(taskData))
	if err != nil {
		log.Fatalf("failed to enqueue task to redis: %v", err)
	}
	fmt.Printf("task '%s' (ID: %s) enqueued to redis\n", newTask.Name, newTask.ID)

	taskQueue := make(chan task.Task, 10)
	pool := queue.NewWorkerPool(3, taskQueue)

	pool.Start()

	go func() {
		for {
			taskID, err := redisClient.Client.LPop(ctx, queueName).Result()
			if err != nil {
				log.Println("no tasks in queue or error fetching task:", err)
				continue
			}

			taskData, err := redisClient.GetTaskData(taskID)
			if err != nil {
				log.Printf("failed to fetch task data for ID %s: %v\n", taskID, err)
				continue
			}

			dataField, exists := taskData["data"]
			if !exists {
				log.Printf("task %s has no data field\n", taskID)
				continue
			}

			var t task.Task
			if err := json.Unmarshal([]byte(dataField), &t); 
			err != nil {
				log.Printf("failed to unmarshal task %s: %v\n", taskID, err)
				continue
			}

			if err := redisClient.UpdateTaskState(t.ID, string(task.Processing)); 
			err != nil {
				log.Printf("failed to update task %s status: %v\n", t.ID, err)
			}

			taskQueue <- t
		}
	}()

	pool.Stop()
}
