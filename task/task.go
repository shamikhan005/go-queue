package task

import (
	"fmt"
	"time"
)

type TaskStatus string

const (
	Pending    TaskStatus = "Pending"
	Processing TaskStatus = "Processing"
	completed  TaskStatus = "Completed"
	Failed     TaskStatus = "Failed"
)

type Task struct {
	ID        string        `json:"id"`
	Name      string     `json:"name"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

func NewTask(name string) *Task {
	return &Task{
		ID:        generateID(),
		Name:      name,
		Status:    Pending,
		CreatedAt: time.Now(),
	}
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
