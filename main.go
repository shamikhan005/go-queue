package main

import (
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

	newTask := task.NewTask(*taskName)

	taskData, err := json.Marshal(newTask)
	if err != nil {
		log.Fatalf("failed to serialize task: %v", err)
	}

	queueName := "task-queue"
	err = redisClient.AddTask(queueName, string(taskData))
	if err != nil {
		log.Fatalf("failed to enqueue task to redis: %v", err)
	}
	fmt.Printf("task '%s' enqueued to redis\n", newTask.Name)

	taskQueue := make(chan task.Task, 10) 
	pool := queue.NewWorkerPool(3, taskQueue)

	pool.Start()

	go func() {
		for {
			taskJSON, err := redisClient.GetTask(queueName)
			if err != nil {
				log.Println("no tasks in queue or error fetching task:", err)
				break
			}
			var t task.Task
			err = json.Unmarshal([]byte(taskJSON), &t)
			if err != nil {
				log.Println("failed to unmarshal task:", err)
				continue
			}
			taskQueue <- t
		}
	}()

	pool.Stop()
}