package queue

import (
	"first/go-queue/task"
	"fmt"
	"time"
)

type WorkerPool struct {
	numWorkers int
	taskQueue  chan task.Task
	quit       chan bool
}

func NewWorkerPool(numWorkers int, taskQueue chan task.Task) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		taskQueue:  taskQueue,
		quit:       make(chan bool),
	}
}

func (p *WorkerPool) Start() {
	for i := 1; i <= p.numWorkers; i++ {
		go p.worker(i)
	}
}

func (p *WorkerPool) Stop() {
	time.Sleep(2 * time.Second)
	close(p.quit)
	fmt.Println("worker pool stopped")
}

func (p *WorkerPool) worker(id int) {
	for {
		select {
		case task := <-p.taskQueue:
			fmt.Printf("worker %d processing task: %s\n", id, task.Name)
			time.Sleep(1 * time.Second)
			fmt.Printf("worker %d completed task: %s\n", id, task.Name)
		case <-p.quit:
			fmt.Printf("worker %d stopping\n", id)
			return
		}
	}
}
