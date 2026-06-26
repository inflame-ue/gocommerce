package carts

import "github.com/inflame-ue/gocommerce/internal/database"

// main cart handler
type CartHandler struct {
	db *database.DB
}

func NewCartHandler(db *database.DB) *CartHandler {
	return &CartHandler{db: db}
}

// database models
type CartItem struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}
