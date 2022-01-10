package database

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/kfsantos/books/models"
)

var DB *gorm.DB

func Setup() {
	host := "localhost"
	port := "5432"
	dbname := "root"
	user := "root"
	password := "root"
	//docker
	// str := "user=root dbname=root password=root host=localhost sslmode=disable"
	// str := "user=postgres dbname=alura_loja password=1234 host=localhost sslmode=disable"
	db, err := gorm.Open("postgres",
		" host="+host+
			" port="+port+
			" user="+user+
			" dbname="+dbname+
			" sslmode=disable password="+password)
	if err != nil {
		log.Fatal(err)
	}
	//Para visualizar consultas SQL
	db.LogMode(false)
	//cria a tabela chamada books
	db.AutoMigrate([]models.Book{})
	DB = db
}
