package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type AddDocumentUseCase struct {
	repo project.Repository
}

func NewAddDocumentUseCase(repo project.Repository) *AddDocumentUseCase {
	return &AddDocumentUseCase{repo: repo}
}

type AddDocumentInput struct {
	ProjectID string
	UserID    string
	Name      string
	Type      string
	Size      string
}

func (uc *AddDocumentUseCase) Execute(ctx context.Context, input AddDocumentInput) (*project.Project, error) {
	if input.ProjectID == "" {
		return nil, errors.New("project ID is required")
	}
	if input.Name == "" || input.Type == "" {
		return nil, errors.New("document name and type are required")
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

	now := time.Now()
	newDoc := project.Document{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Type:      input.Type,
		Size:      input.Size,
		UpdatedAt: now,
	}
	if p.Documents == nil {
		p.Documents = []project.Document{}
	}
	p.Documents = append(p.Documents, newDoc)

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
