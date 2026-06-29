# gocommerce

E-Commerce API built with Go, following the [roadmap.sh project spec](https://roadmap.sh/projects/e-commerce-api).

## Tech Stack

| Layer | Choice |
|---|---|
| Language | Go 1.22+ |
| Router | chi |
| Database | PostgreSQL |
| Driver | pgx |
| Auth | JWT (golang-jwt) |
| Migrations | goose |

## Setup

```bash
# Prerequisites: Go 1.22+, PostgreSQL running

git clone <repo>
cd gocommerce

# Copy environment config
cp .env.example .env  # or create from scratch

# Run migrations
goose -dir migrations postgres "$DATABASE_URL" up

# Start server
go run ./cmd/api/
```

## Environment Variables

| Variable | Description |
|---|---|
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | Secret key for signing JWT tokens |
| `PORT` | Server port (default: 8080) |

## API Endpoints

### Auth (`/auth`)
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/auth/signup` | No | Create account |
| POST | `/auth/login` | No | Login, returns JWT |

### Products (`/products`)
| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/products` | No | List all / search (`?q=term`) |
| GET | `/products/{id}` | No | Get by ID |
| POST | `/products` | Admin | Create product |
| PUT | `/products/{id}` | Admin | Update product |
| DELETE | `/products/{id}` | Admin | Delete product |

### Cart (`/cart`)
| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/cart` | User | View cart items |
| POST | `/cart/{productID}` | User | Add product (increments quantity) |
| DELETE | `/cart/{productID}` | User | Remove product from cart |

### Orders (`/orders`)
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/checkout` | User | Checkout (creates order, clears cart, decrements stock) |
| GET | `/orders` | User | List user's orders |
| GET | `/orders/{id}` | User | Get order with items |
| PATCH | `/orders/{id}/status` | Admin | Update order status |

## Project Structure

```
cmd/api/main.go          — Entry point, wiring
internal/
  auth/                  — Signup, login, JWT, middleware
  carts/                 — Cart CRUD
  database/              — PGX connection
  orders/                — Checkout, order CRUD
  products/              — Product CRUD, search
  response/              — JSON helper
migrations/              — Goose SQL migrations
```

## License

MIT