package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type UpdateDocumentUseCase struct {
	repo project.Repository
}

func NewUpdateDocumentUseCase(repo project.Repository) *UpdateDocumentUseCase {
	return &UpdateDocumentUseCase{repo: repo}
}

type UpdateDocumentInput struct {
	ProjectID string
	UserID    string
	DocumentID string
	Name      *string
	Type      *string
	Size      *string
}

func (uc *UpdateDocumentUseCase) Execute(ctx context.Context, input UpdateDocumentInput) (*project.Project, error) {
	if input.ProjectID == "" || input.DocumentID == "" {
		return nil, errors.New("project ID and document ID are required")
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
	for i := range p.Documents {
		if p.Documents[i].ID == input.DocumentID {
			found = true
			if input.Name != nil {
				p.Documents[i].Name = *input.Name
			}
			if input.Type != nil {
				p.Documents[i].Type = *input.Type
			}
			if input.Size != nil {
				p.Documents[i].Size = *input.Size
			}
			p.Documents[i].UpdatedAt = time.Now()
			break
		}
	}
	if !found {
		return nil, errors.New("document not found")
	}

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
