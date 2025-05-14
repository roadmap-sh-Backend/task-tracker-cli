package main

import "time"

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
