package auth

import "github.com/inflame-ue/gocommerce/internal/database"

// database query models
type userModel struct {
	id            int
	name          string
	email         string
	password_hash string
	is_admin      bool
}

// endpoint related stuff
type signUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
	Msg   string `json:"message"`
}

// main struct of the package
type AuthHandler struct {
	db        *database.DB
	jwtSecret string
}

func NewAuthHandler(db *database.DB, secret string) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: secret}
}
