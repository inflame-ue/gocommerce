package carts

import "github.com/inflame-ue/gocommerce/internal/database"

type CartHandler struct {
	db *database.DB
}

func NewCartHandler(db *database.DB) *CartHandler {
	return &CartHandler{db: db}
}
