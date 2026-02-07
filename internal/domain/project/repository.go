package project

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for project data access
type Repository interface {
	Create(ctx context.Context, project *Project) error
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Project, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]Project, error)
}
