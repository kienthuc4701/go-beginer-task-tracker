package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const filePath = "tasks.json"

func main() {
	fmt.Println("Task Tracker CLI")
	reader := bufio.NewReader(os.Stdin)

	if err := ensureFileExists(); err != nil {
		fmt.Printf("Error initializing file: %v\n", err)
		return
	}

	for {
		fmt.Println("\nEnter a command (add, update, delete, list, exit):")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "add":
			handleAddTask(reader)
		case "update":
			handleUpdateTask(reader)
		case "delete":
			handleDeleteTask(reader)
		case "list":
			handleListTasks(reader)
		case "exit":
			fmt.Println("Exiting Task Tracker CLI. Thank you!")
			return
		default:
			fmt.Println("Unknown command, please try again.")
		}
	}
}

// File Initialization
func ensureFileExists() error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create tasks.json: %v", err)
		}
		defer file.Close()
		_, err = file.WriteString("[]")
		if err != nil {
			return fmt.Errorf("failed to initialize tasks.json: %v", err)
		}
	}
	return nil
}

// Task Handlers
func handleAddTask(reader *bufio.Reader) {
	fmt.Print("Enter task description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	if err := addTask(description); err != nil {
		fmt.Println("Error adding task:", err)
	} else {
		fmt.Println("Task added successfully!")
	}
}

func handleUpdateTask(reader *bufio.Reader) {
	fmt.Print("Enter task ID to update: ")
	id, err := readIntInput(reader)
	if err != nil {
		fmt.Println("Invalid task ID.")
		return
	}

	fmt.Println("1. Update Description")
	fmt.Println("2. Update Status")
	fmt.Print("Choose option: ")

	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	switch option {
	case "1":
		fmt.Print("Enter new task description: ")
		description, _ := reader.ReadString('\n')
		description = strings.TrimSpace(description)
		if err := updateTaskDescription(id, description); err != nil {
			fmt.Println("Error updating description:", err)
		} else {
			fmt.Println("Task description updated successfully!")
		}
	case "2":
		fmt.Print("Enter new status (todo, in-progress, done): ")
		status, _ := reader.ReadString('\n')
		status = strings.TrimSpace(status)
		if err := updateTaskStatus(id, status); err != nil {
			fmt.Println("Error updating status:", err)
		} else {
			fmt.Println("Task status updated successfully!")
		}
	default:
		fmt.Println("Invalid option")
	}
}

func handleDeleteTask(reader *bufio.Reader) {
	fmt.Print("Enter task ID to delete: ")
	id, err := readIntInput(reader)
	if err != nil {
		fmt.Println("Invalid task ID.")
		return
	}

	if err := deleteTask(id); err != nil {
		fmt.Println("Error deleting task:", err)
	} else {
		fmt.Println("Task deleted successfully!")
	}
}

func handleListTasks(reader *bufio.Reader) {
	fmt.Println("Choose listing option:")
	fmt.Println("1. Done")
	fmt.Println("2. In-progress")
	fmt.Println("3. Todo")
	fmt.Println("4. All (default)")
	fmt.Print("Select option (1-4): ")

	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	var status string
	switch option {
	case "1":
		status = "done"
	case "2":
		status = "in-progress"
	case "3":
		status = "todo"
	case "4", "":
		status = ""
	default:
		fmt.Println("Invalid option")
	}

	tasks, err := listTasksByStatus(status)
	if err != nil {
		fmt.Println("Error listing tasks:", err)
		return
	}

	fmt.Println("Tasks:")
	for _, task := range tasks {
		fmt.Printf("ID: %d, Description: %s, Status: %s\n", task.ID, task.Description, task.Status)
	}
}

// CRUD Operations
func listTasksByStatus(status string) ([]Task, error) {
	tasks, err := listTasks()
	if err != nil {
		return nil, err
	}

	if status == "" {
		return tasks, nil // Return all tasks if no status filter
	}

	var filteredTasks []Task
	for _, task := range tasks {
		if task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks, nil
}

func listTasks() ([]Task, error) {
	var tasks []Task
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func addTask(description string) error {
	tasks, err := listTasks()
	if err != nil {
		return err
	}
	id := len(tasks) + 1
	task := Task{
		ID:          id,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, task)
	return saveTasks(tasks)
}

func updateTaskDescription(id int, description string) error {
	return updateTaskField(id, func(task *Task) {
		task.Description = description
		task.UpdatedAt = time.Now()
	})
}

func updateTaskStatus(id int, status string) error {
	if !isValidStatus(status) {
		return errors.New("invalid status")
	}
	return updateTaskField(id, func(task *Task) {
		task.Status = status
		task.UpdatedAt = time.Now()
	})
}

func deleteTask(id int) error {
	tasks, err := listTasks()
	if err != nil {
		return err
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return saveTasks(tasks)
		}
	}
	return errors.New("task not found")
}

func updateTaskField(id int, updateFunc func(*Task)) error {
	tasks, err := listTasks()
	if err != nil {
		return err
	}
	for i, task := range tasks {
		if task.ID == id {
			updateFunc(&tasks[i])
			return saveTasks(tasks)
		}
	}
	return errors.New("task not found")
}

// Helper Functions
func readIntInput(reader *bufio.Reader) (int, error) {
	input, _ := reader.ReadString('\n')
	return strconv.Atoi(strings.TrimSpace(input))
}

func isValidStatus(status string) bool {
	return status == "todo" || status == "in-progress" || status == "done"
}
