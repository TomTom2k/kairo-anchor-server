package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type CreateProjectUseCase struct {
	repo project.Repository
}

func NewCreateProjectUseCase(repo project.Repository) *CreateProjectUseCase {
	return &CreateProjectUseCase{repo: repo}
}

type CreateProjectInput struct {
	UserID      string             `json:"userId"`
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	Status      project.ProjectStatus `json:"status" binding:"required,oneof=active pending completed"`
	Progress    int                `json:"progress" binding:"min=0,max=100"`
	StartDate   time.Time          `json:"startDate" binding:"required"`
	EndDate     *time.Time         `json:"endDate"`
}

func (uc *CreateProjectUseCase) Execute(ctx context.Context, input CreateProjectInput) (*project.Project, error) {
	if input.Name == "" {
		return nil, errors.New("project name is required")
	}

	if input.Progress < 0 || input.Progress > 100 {
		return nil, errors.New("progress must be between 0 and 100")
	}

	if input.EndDate != nil && input.EndDate.Before(input.StartDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	// Parse UserID to UUID
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	p := &project.Project{
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		Status:      input.Status,
		Progress:    input.Progress,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Tasks:       []project.Task{},
		Documents:   []project.Document{},
	}

	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}
