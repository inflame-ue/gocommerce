package carts

import (
	"context"
	"fmt"
)

func (ch *CartHandler) ListCartItems(ctx context.Context, userID int) ([]CartItem, error) {
	var item CartItem
	var items []CartItem

	rows, err := ch.db.Conn.Query(ctx, "SELECT cart_items.quantity, products.id, products.name, products.price FROM cart_items JOIN products ON carts_items.product_id WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("creating the query: %w", err)
	}
	for rows.Next() {
		if err := rows.Scan(&item.Quantity, &item.ProductID, &item.Name, &item.Price); err != nil {
			return nil, fmt.Errorf("scanning one of the rows: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}
