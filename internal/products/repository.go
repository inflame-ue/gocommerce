package products

import (
	"context"
	"fmt"
)

func (ph *ProductHandler) ListProducts(ctx context.Context) ([]productModel, error) {
	rows, err := ph.db.Conn.Query(ctx, "SELECT id, name, price, stock FROM products;")
	if err != nil {
		return nil, fmt.Errorf("sending and initializing query: %w", err)
	}

	var products []productModel
	var product productModel
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
			return nil, fmt.Errorf("scanning the rows: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (ph *ProductHandler) GetProductByID(ctx context.Context, productID int) (*productModel, error) {
	var product productModel
	row := ph.db.Conn.QueryRow(ctx, "SELECT id, name, price, stock FROM products WHERE id = $1", productID)
	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
		return nil, fmt.Errorf("scanning the row: %w", err)
	}
	return &product, nil
}

func (ph *ProductHandler) CreateProduct(ctx context.Context, name string, price float64, stock int) (int, error) {
	var resultID int
	row := ph.db.Conn.QueryRow(ctx, "INSERT INTO products(name, price, stock) VALUES ($1, $2, $3) RETURNING id;", name, price, stock)
	if err := row.Scan(&resultID); err != nil {
		return 0, fmt.Errorf("scanning the row for id: %w", err)
	}
	return resultID, nil
}

func (ph *ProductHandler) UpdateProductByID(ctx context.Context, id int, name string, price float64, stock int) (*productModel, error) {
	var product productModel
	row := ph.db.Conn.QueryRow(ctx, "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4 RETURNING id, name, price, stock",
		name, price, stock, id)
	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
		return nil, fmt.Errorf("scanning the row: %w", err)
	}
	return &product, nil
}

func (ph *ProductHandler) DeleteProductByID(ctx context.Context, id int) (int, error) {
	cmdTag, err := ph.db.Conn.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return int(cmdTag.RowsAffected()), fmt.Errorf("deleting the record: %w", err)
	}
	return int(cmdTag.RowsAffected()), nil
}
