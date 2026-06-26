package carts

import (
	"context"
	"fmt"
)

func (ch *CartHandler) ListCartItems(ctx context.Context, userID int) ([]CartItem, error) {
	var item CartItem
	var items []CartItem

	rows, err := ch.db.Conn.Query(ctx, "SELECT cart_items.quantity, products.id, products.name, products.price FROM cart_items JOIN products ON cart_items.product_id = products.id WHERE cart_items.user_id = $1", userID)
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

func (ch *CartHandler) ListCartItem(ctx context.Context, userID, productID int) (*CartItem, error) {
	var item CartItem

	row := ch.db.Conn.QueryRow(ctx, "SELECT cart_items.quantity, products.id, products.name, products.price FROM cart_items JOIN products ON cart_items.product_id = products.id WHERE cart_items.user_id = $1 AND cart_items.product_id = $2", userID, productID)
	if err := row.Scan(&item.Quantity, &item.ProductID, &item.Name, &item.Price); err != nil {
		return nil, fmt.Errorf("scanning product item: %w", err)
	}

	return &item, nil
}

func (ch *CartHandler) AddCartItem(ctx context.Context, userID, productID int) error {
	_, err := ch.db.Conn.Exec(ctx, "INSERT INTO cart_items(user_id, product_id) VALUES ($1, $2) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = quantity + 1", userID, productID)
	if err != nil {
		return fmt.Errorf("upserting a cart item: %w", err)
	}
	return nil
}

func (ch *CartHandler) DeleteCartItem(ctx context.Context, userID, productID int) error {
	cmdTag, err := ch.db.Conn.Exec(ctx, "DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2", userID, productID)
	if err != nil {
		return fmt.Errorf("executing the query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no cart item available for the given user and product IDs")
	}
	return nil
}
