package postgres

import (
	"context"
	"database/sql"

	"github.com/tomtom2k/kairo-anchor-server/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (email, password) VALUES ($1,$2)",
		u.Email, u.Password,
	)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, password FROM users WHERE email=$1",
		email,
	).Scan(&u.ID, &u.Email, &u.Password)

	if err != nil {
		return nil, nil
	}
	return &u, nil
}