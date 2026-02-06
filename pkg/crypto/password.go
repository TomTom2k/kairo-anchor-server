package crypto

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func (BcryptHasher) Hash(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(b), err
}

func (BcryptHasher) Compare(h, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}
