package user

import "time"

type User struct {
	ID                string
	Email             string
	Password          string
	IsActive          bool
	ActivationToken   *string
	ResetToken        *string
	ResetTokenExpires *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}