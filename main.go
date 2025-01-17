package main

import (
	"first/go-queue/queue"
	"first/go-queue/task"
	"flag"
	"fmt"
	"log"
)

func main() {
	taskName := flag.String("task", "default-task", "name of the task to enqueue")
	flag.Parse()

	newTask := task.NewTask(*taskName)
	if err := task.SaveTask(newTask); err != nil {
		log.Fatalf("failed to save task: %v", err)
	}

	taskQueue := make(chan task.Task, 10) 
	pool := queue.NewWorkerPool(3, taskQueue)

	pool.Start()

	fmt.Printf("enqueuing task: %s\n", newTask.Name)
	taskQueue <- newTask

	pool.Stop()
}