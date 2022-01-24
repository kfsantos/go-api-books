package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kfsantos/books/database"
	"github.com/kfsantos/books/routers"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	database.Setup()
	r := routers.Setup()
	if err := r.Run("127.0.0.1:5000"); err != nil {
		log.Fatal(err)
	}
}
