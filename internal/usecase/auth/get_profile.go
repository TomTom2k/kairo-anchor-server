package auth

import (
	"context"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type GetProfileUseCase struct {
	repo user.Repository
}

func NewGetProfileUseCase(r user.Repository) *GetProfileUseCase {
	return &GetProfileUseCase{r}
}

func (g *GetProfileUseCase) Execute(ctx context.Context, userID string) (*user.User, error) {
	return g.repo.FindByID(ctx, userID)
}
