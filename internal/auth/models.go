package auth

import "github.com/inflame-ue/gocommerce/internal/database"

type AuthHandler struct {
	db *database.DB
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{db: db}
}
