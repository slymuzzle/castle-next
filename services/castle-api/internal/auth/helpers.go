package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Parse jwt
func ParseJwtToken(secret, token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

// Generate jwt
func GenerateJwtToken(secret string, subject string, expiresAt time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   string(subject),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// Hash password
func HashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(pass), err
}

// Compare password and hash
func CompareHashAndPassword(hash string, password string) error {
	byteHash := []byte(hash)
	bytePassword := []byte(password)

	return bcrypt.CompareHashAndPassword(byteHash, bytePassword)
}
