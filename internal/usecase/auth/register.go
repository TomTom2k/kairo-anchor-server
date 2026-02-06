package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type RegisterUseCase struct {
	repo         user.Repository
	hasher       user.PasswordHasher
	emailService user.EmailService
}

func NewRegisterUseCase(r user.Repository, h user.PasswordHasher, e user.EmailService) *RegisterUseCase {
	return &RegisterUseCase{r, h, e}
}

func (r *RegisterUseCase) Execute(ctx context.Context, email, password string) error {
	// Check if user already exists
	existingUser, err := r.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Hash password
	passwordHash, err := r.hasher.Hash(password)
	if err != nil {
		return err
	}

	// Generate activation token
	activationToken, err := generateToken()
	if err != nil {
		return err
	}

	// Create user (inactive by default)
	u := &user.User{
		Email:           email,
		Password:        passwordHash,
		IsActive:        false,
		ActivationToken: &activationToken,
	}

	if err := r.repo.Create(ctx, u); err != nil {
		return err
	}

	// Send activation email
	return r.emailService.SendActivationEmail(email, activationToken)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}