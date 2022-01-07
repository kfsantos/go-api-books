package main

import (
	"log"

	"github.com/kfsantos/books/database"
	"github.com/kfsantos/books/routers"
)

func main() {
	database.Setup()
	r := routers.Setup()
	if err := r.Run("127.0.0.1:5000"); err != nil {
		log.Fatal(err)
	}
}
