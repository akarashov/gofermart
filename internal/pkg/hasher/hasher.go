package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

// Password hash interface
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type bcryptHasher struct{}

// create an instance of password hasher
func NewBcryptHasher() PasswordHasher {
	return &bcryptHasher{}
}

// Hashes the given password using bcrypt.
func (b *bcryptHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Compares the hashed password with the plain password.
func (b *bcryptHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
