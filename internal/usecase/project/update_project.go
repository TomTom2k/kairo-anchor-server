package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type UpdateProjectUseCase struct {
	repo project.Repository
}

func NewUpdateProjectUseCase(repo project.Repository) *UpdateProjectUseCase {
	return &UpdateProjectUseCase{repo: repo}
}

type UpdateProjectInput struct {
	ID          string                `json:"id" binding:"required"`
	UserID      string                `json:"userId"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Status      project.ProjectStatus `json:"status" binding:"omitempty,oneof=active pending completed"`
	Progress    *int                  `json:"progress" binding:"omitempty,min=0,max=100"`
	StartDate   *time.Time            `json:"startDate"`
	EndDate     *time.Time            `json:"endDate"`
	Tasks       []project.Task        `json:"tasks"`
	Documents   []project.Document    `json:"documents"`
}

func (uc *UpdateProjectUseCase) Execute(ctx context.Context, input UpdateProjectInput) (*project.Project, error) {
	if input.ID == "" {
		return nil, errors.New("project ID is required")
	}

	// Parse IDs to UUID
	projectID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, errors.New("invalid project ID format")
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	existingProject, err := uc.repo.FindByID(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}
	if existingProject == nil {
		return nil, errors.New("project not found")
	}

	// Update fields if provided
	if input.Name != "" {
		existingProject.Name = input.Name
	}
	if input.Description != "" {
		existingProject.Description = input.Description
	}
	if input.Status != "" {
		existingProject.Status = input.Status
	}
	if input.Progress != nil {
		if *input.Progress < 0 || *input.Progress > 100 {
			return nil, errors.New("progress must be between 0 and 100")
		}
		existingProject.Progress = *input.Progress
	}
	if input.StartDate != nil {
		existingProject.StartDate = *input.StartDate
	}
	if input.EndDate != nil {
		existingProject.EndDate = input.EndDate
	}
	if input.Tasks != nil {
		existingProject.Tasks = input.Tasks
	}
	if input.Documents != nil {
		existingProject.Documents = input.Documents
	}

	// Validate end date
	if existingProject.EndDate != nil && existingProject.EndDate.Before(existingProject.StartDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	if err := uc.repo.Update(ctx, existingProject); err != nil {
		return nil, err
	}

	return existingProject, nil
}
