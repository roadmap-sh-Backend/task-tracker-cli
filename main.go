package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

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
		taskDescription := strings.Join(taskDescriptionArgs, " ")

		task, err := CreateTask(FILE_NAME, taskDescription)
		if err != nil {
			log.Fatalf("FAILED ADDING TASK: %v", err)
		}

		task.CreatedAt = time.Now()
		log.Printf("Task added successfully (ID: %d)", task.ID)
	case "update":
		taskDescriptionArgs := args[4:]
		taskDescription := strings.Join(taskDescriptionArgs, " ")

		taskId := args[3]
		id := validateId(taskId)

		_, err = GetTaskByID(id)
		if err != nil {
			log.Fatalf("FAILED READING A TASK: %v", err)
		}

		err = UpdateTask(FILE_NAME, id, taskDescription)
		if err != nil {
			log.Fatalf("FAILED UPDATING TASK: %v", err)
		}

		log.Println("Task updated successfully")
	case "delete":
		taskId := args[3]
		id := validateId(taskId)

		_, err = GetTaskByID(id)
		if err != nil {
			log.Fatalf("FAILED READING A TASK: %v", err)
		}

		err = DeleteTask(FILE_NAME, id)
		if err != nil {
			log.Fatalf("TASK DELETED FAILED: %v", err)
		}

		log.Println("Task deleted successfully")
	case "mark-in-progress":
		taskId := args[3]
		id := validateId(taskId)

		_, err = GetTaskByID(id)
		if err != nil {
			log.Fatalf("FAILED READING A TASK: %v", err)
		}

		err = UpdateTaskStatus(FILE_NAME, id, IN_PROGRESS)
		if err != nil {
			log.Fatalf("MARKING TASK AS IN_PROGRESS FAILED: %v", err)
		}

		log.Println("Marking task as in_progress success")
	case "mark-done":
		taskId := args[3]
		id := validateId(taskId)

		_, err = GetTaskByID(id)
		if err != nil {
			log.Fatalf("FAILED READING A TASK: %v", err)
		}

		err = UpdateTaskStatus(FILE_NAME, id, DONE)
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

func validateId(taskId string) int {
	id, err := strconv.ParseInt(taskId, 10, 32)
	if err != nil {
		log.Fatalf("INVALID ID %s", taskId)
	}

	return int(id)
}
