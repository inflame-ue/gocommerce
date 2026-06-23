package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/inflame-ue/gocommerce/internal/auth"
	"github.com/inflame-ue/gocommerce/internal/database"
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

	auth := auth.NewAuthHandler(db)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", auth.HandleSignUp)
		r.Post("/login", auth.HandleLogin)
	})

	port := os.Getenv("PORT")
	log.Printf("listening on port: %s", port)
	err = http.ListenAndServe(":"+port, r)
	log.Fatal(err)
}
