package email

import (
	"fmt"
	"log"
)

// MockEmailService logs emails to console for development
type MockEmailService struct {
	baseURL string
}

func NewMockEmailService(baseURL string) *MockEmailService {
	return &MockEmailService{baseURL: baseURL}
}

func (s *MockEmailService) SendActivationEmail(email, token string) error {
	activationLink := fmt.Sprintf("%s/api/auth/activate?token=%s", s.baseURL, token)
	
	log.Printf("\n=== ACTIVATION EMAIL ===")
	log.Printf("To: %s", email)
	log.Printf("Subject: Activate Your Account")
	log.Printf("Body:")
	log.Printf("  Welcome! Please activate your account by clicking the link below:")
	log.Printf("  %s", activationLink)
	log.Printf("  Or use this token: %s", token)
	log.Printf("========================\n")
	
	return nil
}

func (s *MockEmailService) SendPasswordResetEmail(email, token string) error {
	resetLink := fmt.Sprintf("%s/api/auth/change-password?token=%s", s.baseURL, token)
	
	log.Printf("\n=== PASSWORD RESET EMAIL ===")
	log.Printf("To: %s", email)
	log.Printf("Subject: Reset Your Password")
	log.Printf("Body:")
	log.Printf("  You requested a password reset. Click the link below to reset your password:")
	log.Printf("  %s", resetLink)
	log.Printf("  Or use this token: %s", token)
	log.Printf("  This link will expire in 1 hour.")
	log.Printf("============================\n")
	
	return nil
}
