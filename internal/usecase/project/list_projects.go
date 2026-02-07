package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type ListProjectsUseCase struct {
	repo project.Repository
}

func NewListProjectsUseCase(repo project.Repository) *ListProjectsUseCase {
	return &ListProjectsUseCase{repo: repo}
}

func (uc *ListProjectsUseCase) Execute(ctx context.Context, userID string) ([]project.Project, error) {
	// Parse UserID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	projects, err := uc.repo.FindAllByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	if projects == nil {
		return []project.Project{}, nil
	}

	return projects, nil
}
