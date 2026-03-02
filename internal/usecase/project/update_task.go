package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type UpdateTaskUseCase struct {
	repo project.Repository
}

func NewUpdateTaskUseCase(repo project.Repository) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{repo: repo}
}

type UpdateTaskInput struct {
	ProjectID string
	UserID    string
	TaskID    string
	Title     *string
	Status    *project.TaskStatus
	Priority  *project.TaskPriority
	DueDate   *time.Time
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, input UpdateTaskInput) (*project.Project, error) {
	if input.ProjectID == "" || input.TaskID == "" {
		return nil, errors.New("project ID and task ID are required")
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

	found := false
	for i := range p.Tasks {
		if p.Tasks[i].ID == input.TaskID {
			found = true
			if input.Title != nil {
				p.Tasks[i].Title = *input.Title
			}
			if input.Status != nil {
				p.Tasks[i].Status = *input.Status
			}
			if input.Priority != nil {
				p.Tasks[i].Priority = *input.Priority
			}
			if input.DueDate != nil {
				p.Tasks[i].DueDate = input.DueDate
			}
			break
		}
	}
	if !found {
		return nil, errors.New("task not found")
	}

	p.Progress = project.ProgressFromTasks(p.Tasks)
	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
