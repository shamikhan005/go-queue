package task

import (
	"encoding/json"
	"os"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

var taskCounter int

func NewTask(name string) Task {
	taskCounter++
	return Task{
		ID:        taskCounter,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

func SaveTask(task Task) error {
	file, err := os.OpenFile("tasks.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	taskData, err := json.Marshal(task)
	if err != nil {
		return err
	}

	if _, err = file.Write(append(taskData, '\n')); err != nil {
		return err
	}

	return nil
}
