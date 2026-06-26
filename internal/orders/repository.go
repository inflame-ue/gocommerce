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

func (oh *OrderHandler) CheckoutOrder(ctx context.Context, userID int) (*OrderModel, error) {
	tx, err := oh.db.Conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beggining transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var product productModel
	var products []productModel

	rows, err := tx.Query(ctx, "SELECT cart_items.product_id, products.name, products.price, cart_items.quantity, products.stock FROM cart_items JOIN products ON cart_items.product_id = products.id WHERE cart_items.user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("initializing query for getting the cart items: %w", err)
	}
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.UnitPrice, &product.Quantity, &product.Stock); err != nil {
			return nil, fmt.Errorf("scanning the product rows: %w", err)
		}
		products = append(products, product)
	}

	var total float64
	for _, product := range products {
		total += product.UnitPrice * float64(product.Quantity)
	}

	var order OrderModel
	row := tx.QueryRow(ctx, "INSERT INTO orders(user_id, total) VALUES ($1, $2) RETURNING id, status, total, created_at, updated_at", userID, total)
	if err := row.Scan(&order.ID, &order.Status, &order.Total, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return nil, fmt.Errorf("executing insert into orders: %w", err)
	}

	for _, product := range products {
		_, err := tx.Exec(ctx, "INSERT INTO order_items(order_id, product_id, unit_price, quantity) VALUES ($1, $2, $3, $4)", order.ID, product.ID, product.UnitPrice, product.Quantity)
		if err != nil {
			return nil, fmt.Errorf("executing insert into order_items: %w", err)
		}

		_, err = tx.Exec(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", product.Quantity, product.ID)
		if err != nil {
			return nil, fmt.Errorf("decrementing product stock: %w", err)
		}
	}

	_, err = tx.Exec(ctx, "DELETE FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("deleteting cart items: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commiting the transaction: %w", err)
	}

	return &order, nil
}
