package project

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type DeleteDocumentUseCase struct {
	repo project.Repository
}

func NewDeleteDocumentUseCase(repo project.Repository) *DeleteDocumentUseCase {
	return &DeleteDocumentUseCase{repo: repo}
}

func (uc *DeleteDocumentUseCase) Execute(ctx context.Context, projectIDStr, userIDStr, documentID string) (*project.Project, error) {
	if projectIDStr == "" || documentID == "" {
		return nil, errors.New("project ID and document ID are required")
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return nil, errors.New("invalid project ID format")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	p, err := uc.repo.FindByID(ctx, projectID, userID)
	if err != nil || p == nil {
		return nil, errors.New("project not found")
	}

	newDocs := make([]project.Document, 0, len(p.Documents))
	for _, d := range p.Documents {
		if d.ID != documentID {
			newDocs = append(newDocs, d)
		}
	}
	p.Documents = newDocs

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
