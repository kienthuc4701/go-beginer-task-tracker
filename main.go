package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {
	listTasks()
}

// Task
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const filePath = "tasks.json"

// Database operations
func listTasks() ([]Task, error) {
	var tasks []Task
	// if file not exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return tasks, fmt.Errorf("failed to create tasks.json: %v", err)
		}
		defer file.Close()
		// Write empty JSON array to the file to initialize it
		if _, err := file.WriteString("[]"); err != nil {
			return tasks, fmt.Errorf("failed to initialize tasks.json: %v", err)
		}
	}
	// Read the file if it exists
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON data into tasks slice
	err = json.Unmarshal(bytes, &tasks)
	return tasks, err
}
