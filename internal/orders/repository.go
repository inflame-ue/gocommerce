package orders

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (oh *OrderHandler) ListOrders(ctx context.Context, userID int) ([]OrderModel, error) {
	var order OrderModel
	var orders []OrderModel

	rows, err := oh.db.Conn.Query(ctx, "SELECT id, status, total, created_at, updated_at FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("initializing the query: %w", err)
	}

	for rows.Next() {
		if err := rows.Scan(&order.ID, &order.Status, &order.Total, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning the rows: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (oh *OrderHandler) GetOrderByID(ctx context.Context, userID, orderID int) (*OrderModel, error) {
	var order OrderModel

	row := oh.db.Conn.QueryRow(ctx, "SELECT id, status, total, created_at, updated_at FROM orders WHERE user_id = $1 AND id = $2", userID, orderID)
	if err := row.Scan(&order.ID, &order.Status, &order.Total, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return nil, fmt.Errorf("scanning the row: %w", err)
	}

	var product productModel
	rows, err := oh.db.Conn.Query(ctx, "SELECT products.id, products.name, order_items.unit_price, order_items.quantity FROM order_items JOIN products ON order_items.product_id = products.id WHERE order_items.order_id = $1", order.ID)
	if err != nil {
		return nil, fmt.Errorf("initializing the products query: %w", err)
	}
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.UnitPrice, &product.Quantity); err != nil {
			return nil, fmt.Errorf("scanning the rows: %w", err)
		}
		order.Products = append(order.Products, product)
	}

	return &order, nil
}

func (oh *OrderHandler) UpdateOrderStatusByID(ctx context.Context, userID, orderID int, status string) error {
	res, err := oh.db.Conn.Exec(ctx, "UPDATE orders SET status = $1 WHERE user_id = $2 AND id = $3", status, userID, orderID)
	if err != nil {
		return fmt.Errorf("executing the query: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected %w", pgx.ErrNoRows)
	}
	return nil
}

func (oh *OrderHandler) CheckoutOrder(ctx context.Context, userID int) error {
	return nil
}
