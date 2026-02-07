package project

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type GetProjectUseCase struct {
	repo project.Repository
}

func NewGetProjectUseCase(repo project.Repository) *GetProjectUseCase {
	return &GetProjectUseCase{repo: repo}
}

func (uc *GetProjectUseCase) Execute(ctx context.Context, id string, userID string) (*project.Project, error) {
	if id == "" {
		return nil, errors.New("project ID is required")
	}

	// Parse IDs to UUID
	projectID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid project ID format")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	p, err := uc.repo.FindByID(ctx, projectID, userUUID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("project not found")
	}

	return p, nil
}
