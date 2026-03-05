package auth

import "golang.org/x/crypto/bcrypt"

// PasswordManager handles password hashing and verification
type PasswordManager struct{}

func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}

// HashPassword hashes a plain text password using bcrypt
func (m *PasswordManager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// VerifyPassword checks if a plain text password matches a hash
func (m *PasswordManager) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
