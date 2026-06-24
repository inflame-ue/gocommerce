package products

import "github.com/inflame-ue/gocommerce/internal/database"

type ProductHandler struct {
	db *database.DB
}

func NewProductHandler(db *database.DB) *ProductHandler {
	return &ProductHandler{db: db}
}
