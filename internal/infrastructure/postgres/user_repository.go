package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (email, password, is_active, activation_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		u.Email, u.Password, u.IsActive, u.ActivationToken,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
		SET email = $1, password = $2, is_active = $3, activation_token = $4,
		    reset_token = $5, reset_token_expires = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		u.Email, u.Password, u.IsActive, u.ActivationToken,
		u.ResetToken, u.ResetTokenExpires, u.ID,
	).Scan(&u.UpdatedAt)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	query := `
		SELECT id, email, password, is_active, activation_token,
		       reset_token, reset_token_expires, created_at, updated_at
		FROM users WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Password, &u.IsActive, &u.ActivationToken,
		&u.ResetToken, &u.ResetTokenExpires, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	query := `
		SELECT id, email, password, is_active, activation_token,
		       reset_token, reset_token_expires, created_at, updated_at
		FROM users WHERE email = $1
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.Password, &u.IsActive, &u.ActivationToken,
		&u.ResetToken, &u.ResetTokenExpires, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByActivationToken(ctx context.Context, token string) (*user.User, error) {
	var u user.User
	query := `
		SELECT id, email, password, is_active, activation_token,
		       reset_token, reset_token_expires, created_at, updated_at
		FROM users WHERE activation_token = $1
	`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&u.ID, &u.Email, &u.Password, &u.IsActive, &u.ActivationToken,
		&u.ResetToken, &u.ResetTokenExpires, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid activation token")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByResetToken(ctx context.Context, token string) (*user.User, error) {
	var u user.User
	query := `
		SELECT id, email, password, is_active, activation_token,
		       reset_token, reset_token_expires, created_at, updated_at
		FROM users WHERE reset_token = $1
	`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&u.ID, &u.Email, &u.Password, &u.IsActive, &u.ActivationToken,
		&u.ResetToken, &u.ResetTokenExpires, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid reset token")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}