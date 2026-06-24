package auth

import (
	"context"
	"fmt"
)

func (au *AuthHandler) CreateUser(ctx context.Context, name, email, password_hash string) (int, bool, error) {
	var userID int
	var is_admin bool
	row := au.db.Conn.QueryRow(ctx, "INSERT INTO users(name, email, password_hash) VALUES($1, $2, $3) RETURNING id, is_admin;",
		name, email, password_hash)
	if err := row.Scan(&userID, &is_admin); err != nil {
		return 0, false, fmt.Errorf("inserting user into database: %w", err)
	}
	return userID, is_admin, nil
}

func (au *AuthHandler) GetUserByEmail(ctx context.Context, email string) (*userModel, error) {
	row := au.db.Conn.QueryRow(ctx, "SELECT id, name, email, password_hash, is_admin FROM users WHERE email = $1", email)

	var user userModel
	if err := row.Scan(&user.id, &user.name, &user.email, &user.password_hash, &user.is_admin); err != nil {
		return nil, fmt.Errorf("scanning the row into userObject: %w", err)
	}

	return &user, nil
}
