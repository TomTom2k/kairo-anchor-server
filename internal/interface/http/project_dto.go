package http

import "time"

// Project DTOs

type TaskDTO struct {
	ID       string    `json:"id" example:"t1"`
	Title    string    `json:"title" binding:"required" example:"Thiết kế Database"`
	Status   string    `json:"status" binding:"required,oneof=todo in-progress completed" example:"completed"`
	Priority string    `json:"priority" binding:"required,oneof=low medium high" example:"high"`
	DueDate  *time.Time `json:"dueDate,omitempty" example:"2024-02-01T00:00:00Z"`
}

type DocumentDTO struct {
	ID        string    `json:"id" example:"d1"`
	Name      string    `json:"name" binding:"required" example:"Spec tài liệu.pdf"`
	Type      string    `json:"type" binding:"required" example:"pdf"`
	Size      string    `json:"size" example:"2.4 MB"`
	UpdatedAt time.Time `json:"updatedAt" example:"2024-01-20T00:00:00Z"`
}

type CreateProjectRequest struct {
	Name        string    `json:"name" binding:"required" example:"Hệ thống quản lý kho"`
	Description string    `json:"description" example:"Xây dựng hệ thống quản lý kho thông minh"`
	Status      string    `json:"status" binding:"required,oneof=active pending completed" example:"active"`
	Progress    int       `json:"progress" binding:"min=0,max=100" example:"65"`
	StartDate   time.Time `json:"startDate" binding:"required" example:"2024-01-15T00:00:00Z"`
	EndDate     *time.Time `json:"endDate,omitempty" example:"2024-06-30T00:00:00Z"`
}

type UpdateProjectRequest struct {
	Name        string     `json:"name,omitempty" example:"Hệ thống quản lý kho"`
	Description string     `json:"description,omitempty" example:"Xây dựng hệ thống quản lý kho thông minh"`
	Status      string     `json:"status,omitempty" binding:"omitempty,oneof=active pending completed" example:"active"`
	StartDate   *time.Time `json:"startDate,omitempty" example:"2024-01-15T00:00:00Z"`
	EndDate     *time.Time `json:"endDate,omitempty" example:"2024-06-30T00:00:00Z"`
}

type CreateTaskRequest struct {
	Title    string     `json:"title" binding:"required" example:"Thiết kế Database"`
	Status   string     `json:"status" binding:"required,oneof=todo in-progress completed" example:"todo"`
	Priority string     `json:"priority" binding:"required,oneof=low medium high" example:"medium"`
	DueDate  *time.Time `json:"dueDate,omitempty"`
}

type UpdateTaskRequest struct {
	Title    *string    `json:"title,omitempty"`
	Status   *string    `json:"status,omitempty" binding:"omitempty,oneof=todo in-progress completed"`
	Priority *string    `json:"priority,omitempty" binding:"omitempty,oneof=low medium high"`
	DueDate  *time.Time `json:"dueDate,omitempty"`
}

type ReorderTasksRequest struct {
	TaskIDs []string `json:"taskIds" binding:"required" example:"id1,id2,id3"`
}

type CreateDocumentRequest struct {
	Name string `json:"name" binding:"required" example:"Spec.pdf"`
	Type string `json:"type" binding:"required" example:"pdf"`
	Size string `json:"size" example:"2.4 MB"`
}

type UpdateDocumentRequest struct {
	Name      *string    `json:"name,omitempty"`
	Type      *string    `json:"type,omitempty"`
	Size      *string    `json:"size,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type ProjectResponse struct {
	ID          string        `json:"id" example:"1"`
	UserID      string        `json:"userId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string        `json:"name" example:"Hệ thống quản lý kho"`
	Description string        `json:"description" example:"Xây dựng hệ thống quản lý kho thông minh"`
	Status      string        `json:"status" example:"active"`
	Progress    int           `json:"progress" example:"65"`
	StartDate   time.Time     `json:"startDate" example:"2024-01-15T00:00:00Z"`
	EndDate     *time.Time    `json:"endDate,omitempty" example:"2024-06-30T00:00:00Z"`
	Tasks       []TaskDTO     `json:"tasks"`
	Documents   []DocumentDTO `json:"documents"`
	CreatedAt   time.Time     `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time     `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
}
