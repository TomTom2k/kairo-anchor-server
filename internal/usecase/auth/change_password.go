package auth

import (
	"context"
	"errors"
	"time"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type ChangePasswordUseCase struct {
	repo   user.Repository
	hasher user.PasswordHasher
}

func NewChangePasswordUseCase(r user.Repository, h user.PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{r, h}
}

func (c *ChangePasswordUseCase) Execute(ctx context.Context, token, newPassword string) error {
	// Find user by reset token
	u, err := c.repo.FindByResetToken(ctx, token)
	if err != nil {
		return err
	}

	// Check if token is expired
	if u.ResetTokenExpires != nil && u.ResetTokenExpires.Before(time.Now()) {
		return errors.New("reset token has expired")
	}

	// Hash new password
	passwordHash, err := c.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	// Update password and clear reset token
	u.Password = passwordHash
	u.ResetToken = nil
	u.ResetTokenExpires = nil

	return c.repo.Update(ctx, u)
}
