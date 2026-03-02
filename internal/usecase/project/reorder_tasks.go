package project

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type ReorderTasksUseCase struct {
	repo project.Repository
}

func NewReorderTasksUseCase(repo project.Repository) *ReorderTasksUseCase {
	return &ReorderTasksUseCase{repo: repo}
}

// Execute reorders project tasks to match the given taskIds order.
// Tasks not in taskIds are appended at the end in their current order.
func (uc *ReorderTasksUseCase) Execute(ctx context.Context, projectIDStr, userIDStr string, taskIds []string) (*project.Project, error) {
	if projectIDStr == "" {
		return nil, errors.New("project ID is required")
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

	// Build id -> task map
	byID := make(map[string]project.Task)
	for _, t := range p.Tasks {
		byID[t.ID] = t
	}

	// Build new order: first by taskIds, then any remaining in current order
	seen := make(map[string]bool)
	var newTasks []project.Task
	for _, id := range taskIds {
		if t, ok := byID[id]; ok {
			newTasks = append(newTasks, t)
			seen[id] = true
		}
	}
	for _, t := range p.Tasks {
		if !seen[t.ID] {
			newTasks = append(newTasks, t)
		}
	}

	p.Tasks = newTasks
	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
