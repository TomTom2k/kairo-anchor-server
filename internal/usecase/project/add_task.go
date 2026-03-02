package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type AddTaskUseCase struct {
	repo project.Repository
}

func NewAddTaskUseCase(repo project.Repository) *AddTaskUseCase {
	return &AddTaskUseCase{repo: repo}
}

type AddTaskInput struct {
	ProjectID string
	UserID    string
	Title     string
	Status    project.TaskStatus
	Priority  project.TaskPriority
	DueDate   *time.Time
}

func (uc *AddTaskUseCase) Execute(ctx context.Context, input AddTaskInput) (*project.Project, error) {
	if input.ProjectID == "" {
		return nil, errors.New("project ID is required")
	}
	if input.Title == "" {
		return nil, errors.New("task title is required")
	}

	projectID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, errors.New("invalid project ID format")
	}
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	p, err := uc.repo.FindByID(ctx, projectID, userID)
	if err != nil || p == nil {
		return nil, errors.New("project not found")
	}

	newTask := project.Task{
		ID:       uuid.New().String(),
		Title:    input.Title,
		Status:   input.Status,
		Priority: input.Priority,
		DueDate:  input.DueDate,
	}
	if p.Tasks == nil {
		p.Tasks = []project.Task{}
	}
	p.Tasks = append(p.Tasks, newTask)
	p.Progress = project.ProgressFromTasks(p.Tasks)

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
