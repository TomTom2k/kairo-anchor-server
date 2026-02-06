package http

import "time"

// Request DTOs
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type ActivateAccountRequest struct {
	Token string `json:"token" binding:"required" example:"activation-token-here"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

type ChangePasswordRequest struct {
	Token       string `json:"token" binding:"required" example:"reset-token-here"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newpassword123"`
}

// Response DTOs
type MessageResponse struct {
	Message string `json:"message" example:"Success"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

type LoginResponse struct {
	Token   string          `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User    ProfileResponse `json:"user"`
	Message string          `json:"message" example:"Login successful"`
}

type ProfileResponse struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string    `json:"email" example:"user@example.com"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
