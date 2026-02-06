package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type ForgotPasswordUseCase struct {
	repo         user.Repository
	emailService user.EmailService
}

func NewForgotPasswordUseCase(r user.Repository, e user.EmailService) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{r, e}
}

func (f *ForgotPasswordUseCase) Execute(ctx context.Context, email string) error {
	// Find user by email
	u, err := f.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("user not found")
	}

	// Generate reset token
	resetToken, err := generateResetToken()
	if err != nil {
		return err
	}

	// Set reset token with 1 hour expiration
	expiresAt := time.Now().Add(1 * time.Hour)
	u.ResetToken = &resetToken
	u.ResetTokenExpires = &expiresAt

	if err := f.repo.Update(ctx, u); err != nil {
		return err
	}

	// Send password reset email
	return f.emailService.SendPasswordResetEmail(email, resetToken)
}

func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
