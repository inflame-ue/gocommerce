package products

import "github.com/inflame-ue/gocommerce/internal/database"

// main handler of the package
type ProductHandler struct {
	db *database.DB
}

func NewProductHandler(db *database.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// database models
type productModel struct {
	ID    int     `json:"id,omitempty"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
