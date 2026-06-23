-- +goose Up
CREATE TYPE order_status AS ENUM ('pending', 'paid',  'shipped', 'cancelled');

CREATE TABLE orders(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status order_status NOT NULL DEFAULT 'pending',
    total NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE order_items(
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    unit_price NUMERIC(10, 2) NOT NULL,
    quantity INTEGER DEFAULT 1,
    UNIQUE(order_id, product_id)
);

-- +goose Down
DROP TABLE order_items;
DROP TABLE orders;
DROP TYPE order_status;
