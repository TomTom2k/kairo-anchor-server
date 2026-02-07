package project

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type DeleteProjectUseCase struct {
	repo project.Repository
}

func NewDeleteProjectUseCase(repo project.Repository) *DeleteProjectUseCase {
	return &DeleteProjectUseCase{repo: repo}
}

func (uc *DeleteProjectUseCase) Execute(ctx context.Context, id string, userID string) error {
	if id == "" {
		return errors.New("project ID is required")
	}

	// Parse IDs to UUID
	projectID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid project ID format")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	return uc.repo.Delete(ctx, projectID, userUUID)
}
