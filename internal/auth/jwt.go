package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(email, secret string, is_admin bool, userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":      userID,
			"email":    email,
			"is_admin": is_admin,
			"iat":      time.Now().UTC(),
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("signing the token: %w", err)
	}
	return tokenString, nil
}

func validateJWT(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) { return []byte(secret), nil })
	if err != nil {
		return nil, fmt.Errorf("parsing token string: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("JWT token is invalid")
	}
	return token.Claims.(jwt.MapClaims), nil
}
