package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

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

	err = WriteTask(fileName, tasks)
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

	return WriteTask(fileName, tasks)
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

	return WriteTask(fileName, tasks)
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

	return WriteTask(fileName, newTasks)
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

func WriteTask(fileName string, tasks *Tasks) error {
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
