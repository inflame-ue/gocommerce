package main

import (
	"context"
	"log"
	"os"

	"github.com/inflame-ue/gocommerce/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading .env variable: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	_, err = database.NewDatabase(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("connected to the database succesfully")
}
