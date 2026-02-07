package project

import (
	"time"

	"github.com/google/uuid"
)

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	StatusActive    ProjectStatus = "active"
	StatusPending   ProjectStatus = "pending"
	StatusCompleted ProjectStatus = "completed"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusCompleted  TaskStatus = "completed"
)

// TaskPriority represents the priority level of a task
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

// Task represents a single task within a project
type Task struct {
	ID       string       `json:"id"`
	Title    string       `json:"title"`
	Status   TaskStatus   `json:"status"`
	Priority TaskPriority `json:"priority"`
	DueDate  *time.Time   `json:"dueDate,omitempty"`
}

// Document represents a document associated with a project
type Document struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Size      string    `json:"size"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Project represents a project entity
type Project struct {
	ID          uuid.UUID     `json:"id"`
	UserID      uuid.UUID     `json:"userId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      ProjectStatus `json:"status"`
	Progress    int           `json:"progress"`
	StartDate   time.Time     `json:"startDate"`
	EndDate     *time.Time    `json:"endDate,omitempty"`
	Tasks       []Task        `json:"tasks"`
	Documents   []Document    `json:"documents"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}
