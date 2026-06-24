package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const claimsKey contextKey = "claims"

func parseJWTToken(r *http.Request) (string, error) {
	headerValue := r.Header.Get("Authorization")

	headerParts := strings.Split(headerValue, " ")
	if len(headerParts) < 2 {
		return "", errors.New("malformed authorization header")
	}

	token := strings.TrimSpace(headerParts[1])
	return token, nil
}

func ClaimsFromContext(ctx context.Context) jwt.MapClaims {
	return ctx.Value(claimsKey).(jwt.MapClaims)
}

func (ah *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := parseJWTToken(r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		claims, err := ah.validateJWT(jwtToken)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "failed to verify the JWT"})
			return
		}
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
