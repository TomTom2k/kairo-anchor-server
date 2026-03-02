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

// UpdateProjectInput chỉ chứa thông tin cơ bản; task/document dùng API riêng
type UpdateProjectInput struct {
	ID          string                `json:"id" binding:"required"`
	UserID      string                `json:"userId"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Status      project.ProjectStatus `json:"status" binding:"omitempty,oneof=active pending completed"`
	StartDate   *time.Time            `json:"startDate"`
	EndDate     *time.Time            `json:"endDate"`
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
	// Progress luôn tính từ task hoàn thành, không nhập tay
	if input.StartDate != nil {
		existingProject.StartDate = *input.StartDate
	}
	if input.EndDate != nil {
		existingProject.EndDate = input.EndDate
	}

	// Progress giữ nguyên (chỉ thay đổi qua API task)

	// Validate end date
	if existingProject.EndDate != nil && existingProject.EndDate.Before(existingProject.StartDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	if err := uc.repo.Update(ctx, existingProject); err != nil {
		return nil, err
	}

	return existingProject, nil
}
