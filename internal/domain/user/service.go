package user

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type TokenService interface {
	Generate(userID string) (string, error)
	Validate(token string) (string, error)
}

type EmailService interface {
	SendActivationEmail(email, token string) error
	SendPasswordResetEmail(email, token string) error
}