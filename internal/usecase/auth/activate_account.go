package auth

import (
	"context"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type ActivateAccountUseCase struct {
	repo user.Repository
}

func NewActivateAccountUseCase(r user.Repository) *ActivateAccountUseCase {
	return &ActivateAccountUseCase{r}
}

func (a *ActivateAccountUseCase) Execute(ctx context.Context, token string) error {
	// Find user by activation token
	u, err := a.repo.FindByActivationToken(ctx, token)
	if err != nil {
		return err
	}

	// Mark user as active and clear activation token
	u.IsActive = true
	u.ActivationToken = nil

	return a.repo.Update(ctx, u)
}
