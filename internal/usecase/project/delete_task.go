package project

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type DeleteTaskUseCase struct {
	repo project.Repository
}

func NewDeleteTaskUseCase(repo project.Repository) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{repo: repo}
}

func (uc *DeleteTaskUseCase) Execute(ctx context.Context, projectIDStr, userIDStr, taskID string) (*project.Project, error) {
	if projectIDStr == "" || taskID == "" {
		return nil, errors.New("project ID and task ID are required")
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

	newTasks := make([]project.Task, 0, len(p.Tasks))
	for _, t := range p.Tasks {
		if t.ID != taskID {
			newTasks = append(newTasks, t)
		}
	}
	p.Tasks = newTasks
	p.Progress = project.ProgressFromTasks(p.Tasks)

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
