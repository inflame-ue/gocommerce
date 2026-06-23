package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func generateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generating hash from password: %w", err)
	}
	return string(hash), nil
}

func comparePasswordAndHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
