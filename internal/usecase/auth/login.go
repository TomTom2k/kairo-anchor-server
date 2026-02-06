package auth

import (
	"context"
	"errors"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type LoginUseCase struct {
	repo         user.Repository
	hasher       user.PasswordHasher
	tokenService user.TokenService
}

func NewLoginUseCase(r user.Repository, h user.PasswordHasher, t user.TokenService) *LoginUseCase {
	return &LoginUseCase{r, h, t}
}

type LoginResult struct {
	Token string
	User  *user.User
}

func (l *LoginUseCase) Execute(ctx context.Context, email, password string) (*LoginResult, error) {
	// Find user by email
	u, err := l.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if !l.hasher.Compare(u.Password, password) {
		return nil, errors.New("invalid email or password")
	}

	// Check if account is active
	if !u.IsActive {
		return nil, errors.New("account not activated, please check your email")
	}

	// Generate JWT token
	token, err := l.tokenService.Generate(u.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
		User:  u,
	}, nil
}
