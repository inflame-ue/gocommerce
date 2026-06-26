package orders

import (
	"time"

	"github.com/inflame-ue/gocommerce/internal/database"
)

// handler struct
type OrderHandler struct {
	db *database.DB
}

func NewOrderHandler(db *database.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

// database models
type productModel struct {
	ID        int     `json:"id,omitempty"`
	Name      string  `json:"name"`
	UnitPrice float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

type OrderModel struct {
	ID        int            `json:"id"`
	Status    string         `json:"status"`
	Total     float64        `json:"total"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Products  []productModel `json:"products,omitempty"`
}
