package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type TaskStatus string

var (
	ALL         TaskStatus = "ALL"
	DONE        TaskStatus = "done"
	IN_PROGRESS TaskStatus = "in-progress"
	TODO        TaskStatus = "todo"
)

type Tasks struct {
	Task []Task `json:"task"`
}

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

const FILE_NAME = "tasks.json"

func main() {
	err := initializeStorage(FILE_NAME)
	if err != nil {
		log.Fatalf("ERROR INITIALIZE STORAGE")
	}

	args := os.Args

	if len(args) < 2 {
		log.Fatalf("LACKS COMMAND")
	}

	if args[1] != "task-cli" {
		log.Fatalf("WRONG FIRST ARGUMENT. FIRST ARG MUST 'task-cli'")
	}

	command := args[2]

	switch command {
	case "add":
		taskDescriptionArgs := args[3:]
		log.Printf("adding a new taskDescription: [%s]", taskDescriptionArgs)
		taskDescription := strings.Join(taskDescriptionArgs, " ")

		task, err := CreateTask(FILE_NAME, taskDescription)
		if err != nil {
			log.Fatalf("FAILED ADDING TASK: %v", err)
		}

		task.CreatedAt = time.Now()
		log.Printf("Task added successfully (ID: %d)", task.ID)
	case "update":
		taskDescriptionArgs := args[4:]
		log.Printf("update a taskDescription: [%s]", taskDescriptionArgs)
		taskDescription := strings.Join(taskDescriptionArgs, " ")

		taskId := args[3]

		id, err := strconv.ParseInt(taskId, 10, 32)
		if err != nil {
			log.Fatalf("INVALID ID %s", taskId)
		}

		_, err = GetTaskByID(int(id))
		if err != nil {
			log.Fatalf("FAILED READING A TASK: %v", err)
		}

		err = UpdateTask(FILE_NAME, int(id), taskDescription)
		if err != nil {
			log.Fatalf("FAILED UPDATING TASK: %v", err)
		}

		log.Println("Task updated successfully")
	case "delete":
		taskId := args[3]
		id, err := strconv.ParseInt(taskId, 10, 32)
		if err != nil {
			log.Fatalf("INVALID ID %s", taskId)
		}

		_, err = GetTaskByID(int(id))
		if err != nil {
			log.Fatalf("FAILED READING A TASK: %v", err)
		}

		err = DeleteTask(FILE_NAME, int(id))
		if err != nil {
			log.Fatalf("TASK DELETED FAILED: %v", err)
		}

		log.Println("Task deleted successfully")
	case "mark-in-progress":
		taskId := args[3]
		id, err := strconv.ParseInt(taskId, 10, 32)
		if err != nil {
			log.Fatalf("INVALID ID %s", taskId)
		}

		_, err = GetTaskByID(int(id))
		if err != nil {
			log.Fatalf("FAILED READING A TASK: %v", err)
		}

		err = UpdateTaskStatus(FILE_NAME, int(id), IN_PROGRESS)
		if err != nil {
			log.Fatalf("MARKING TASK AS IN_PROGRESS FAILED: %v", err)
		}

		log.Println("Marking task as in_progress success")
	case "mark-done":
		taskId := args[3]
		id, err := strconv.ParseInt(taskId, 10, 32)
		if err != nil {
			log.Fatalf("INVALID ID %s", taskId)
		}

		_, err = GetTaskByID(int(id))
		if err != nil {
			log.Fatalf("FAILED READING A TASK: %v", err)
		}

		err = UpdateTaskStatus(FILE_NAME, int(id), DONE)
		if err != nil {
			log.Fatalf("MARKING TASK AS DONE FAILED: %v", err)
		}

		log.Println("Marking task as done success")
	case "list":
		status := ALL
		if len(args) > 3 {
			status = TaskStatus(args[3])
			log.Printf("listing the tasks with status: [%s]", status)
		}

		tasks, err := GetTasks(FILE_NAME, status)
		if err != nil {
			log.Fatalf("FAILED READING TASKS: %v", err)
		}

		log.Println("LIST TASKS", tasks)
	default:
		log.Fatalf("UNKNOWN COMMAND")
	}
}

func initializeStorage(name string) error {
	_, err := os.Stat(name)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		return os.WriteFile(name, []byte(`{"task": []}`), 0644)
	}

	return err
}

func CreateTask(fileName, taskDescription string) (*Task, error) {
	tasks, err := GetTasks(fileName, ALL)
	if err != nil {
		return nil, err
	}

	task := Task{
		ID:          len(tasks.Task) + 1,
		Description: taskDescription,
		Status:      TODO,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks.Task = append(tasks.Task, task)

	err = WriteData(fileName, tasks)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(fileName string, id int, taskDescription string) error {
	tasks, err := GetTasks(fileName, ALL)
	if err != nil {
		return err
	}

	for i := range tasks.Task {
		if tasks.Task[i].ID == id {
			tasks.Task[i].Description = taskDescription
			tasks.Task[i].UpdatedAt = time.Now()
			break
		}
	}

	return WriteData(fileName, tasks)
}

func UpdateTaskStatus(fileName string, id int, taskStatus TaskStatus) error {
	tasks, err := GetTasks(fileName, ALL)
	if err != nil {
		return err
	}

	for i := range tasks.Task {
		if tasks.Task[i].ID == id {
			tasks.Task[i].Status = taskStatus
			tasks.Task[i].UpdatedAt = time.Now()
			break
		}
	}

	return WriteData(fileName, tasks)
}

func DeleteTask(fileName string, id int) error {
	tasks, err := GetTasks(fileName, ALL)
	if err != nil {
		return err
	}

	newTasks := &Tasks{
		Task: []Task{},
	}

	for _, task := range tasks.Task {
		if task.ID != id {
			newTasks.Task = append(newTasks.Task, task)
		}
	}

	return WriteData(fileName, newTasks)
}

func GetTasks(fileName string, status TaskStatus) (*Tasks, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("FAILED READING A FILE", err)
		return nil, err
	}

	var tasks Tasks
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&tasks)
	if err != nil && err != io.EOF {
		log.Println("FAILED DECODE TASKS", err)
		return nil, err
	}

	if status == ALL || status == "" {
		return &tasks, nil
	}

	var filteredTasks Tasks

	for _, task := range tasks.Task {
		if task.Status == status {
			filteredTasks.Task = append(filteredTasks.Task, task)
		}
	}

	return &filteredTasks, nil
}

func GetTaskByID(id int) (*Task, error) {
	tasks, err := GetTasks(FILE_NAME, ALL)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks.Task {
		if task.ID == id {
			return &task, nil
		}
	}

	return nil, fmt.Errorf("TASK DID NOT FOUND")
}

func WriteData(fileName string, tasks *Tasks) error {
	jsonData, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		log.Println("Error marshalling", err)
		return err
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Println("Error writing into a file", err)
		return err
	}

	return nil
}
