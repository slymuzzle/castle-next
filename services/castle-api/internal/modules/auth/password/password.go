package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash password
func Hash(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(pass), err
}

// Compare password and hash
func Compare(hash string, password string) error {
	byteHash := []byte(hash)
	bytePassword := []byte(password)

	return bcrypt.CompareHashAndPassword(byteHash, bytePassword)
}
