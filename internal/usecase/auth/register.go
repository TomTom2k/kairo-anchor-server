package auth

import (
	"context"
	"errors"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type RegisterUseCase struct {
	repo   user.Repository
	hasher user.PasswordHasher
}

func NewRegisterUseCase(r user.Repository, h user.PasswordHasher) *RegisterUseCase {
	return &RegisterUseCase{r, h}
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

	// Create user
	user := &user.User{
		Email:    email,
		Password: passwordHash,
	}

	return r.repo.Create(ctx, user)
}