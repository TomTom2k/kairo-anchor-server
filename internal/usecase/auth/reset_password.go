package auth

import (
	"context"
	"errors"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type ResetPasswordUseCase struct {
	userRepo user.Repository
	hasher   user.PasswordHasher
}

func NewResetPasswordUseCase(r user.Repository, h user.PasswordHasher) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		userRepo: r,
		hasher:   h,
	}
}


func (uc *ResetPasswordUseCase) Execute(ctx context.Context, userID string, oldPassword, newPassword string) error {
	// Find user
	u, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !uc.hasher.Compare(u.Password, oldPassword) {
		return errors.New("invalid old password")
	}


	// Hash new password
	hashedPassword, err := uc.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	// Update password
	u.Password = hashedPassword
	return uc.userRepo.Update(ctx, u)
}
