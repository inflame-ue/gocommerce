package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Conn *pgx.Conn
}

func NewDatabase(ctx context.Context, url string) (*DB, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("connecting to the database at %s: %v", url, err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pinging the database at %s: %v", url, err)
	}

	return &DB{Conn: conn}, nil
}
