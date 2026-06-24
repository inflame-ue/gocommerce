package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const claimsKey contextKey = "claims"

func parseJWTToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	token = strings.Split(token, " ")[1]
	return strings.TrimSpace(token)
}

func (ah *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken := parseJWTToken(r)
		claims, err := ah.validateJWT(jwtToken)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to verify the JWT"})
			return
		}
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
