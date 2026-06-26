package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/inflame-ue/gocommerce/internal/auth"
	"github.com/inflame-ue/gocommerce/internal/carts"
	"github.com/inflame-ue/gocommerce/internal/database"
	"github.com/inflame-ue/gocommerce/internal/orders"
	"github.com/inflame-ue/gocommerce/internal/products"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading .env variable: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	db, err := database.NewDatabase(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	jwtSecret := os.Getenv("JWT_SECRET")
	auth := auth.NewAuthHandler(db, jwtSecret)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", auth.HandleSignUp)
		r.Post("/login", auth.HandleLogin)
	})

	product := products.NewProductHandler(db)
	r.Get("/products", product.HandleGetProducts)
	r.Get("/products/{productID}", product.HandleGetProduct)
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		// require admin rights as well
		r.Post("/products", product.HandleCreateProduct)
		r.Put("/products/{productID}", product.HandleUpdateProduct)
		r.Delete("/products/{productID}", product.HandleDeleteProduct)
	})

	cart := carts.NewCartHandler(db)
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/cart", cart.HandleGetCart)
		r.Post("/cart/{productID}", cart.HandleAddProductToCart)
		r.Delete("/cart/{productID}", cart.HandleDeleteProductFromCart)
	})

	order := orders.NewOrderHandler(db)
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/checkout", order.HandleCheckout)
		r.Get("/orders", order.HandleGetOrders)
		r.Get("/orders/{orderID}", order.HandleGetOrder)
		r.Patch("/orders/{orderID}/status", order.HandleUpdateOrderStatus) // admin
	})

	port := os.Getenv("PORT")
	log.Printf("listening on port: %s", port)
	err = http.ListenAndServe(":"+port, r)
	log.Fatal(err)
}
