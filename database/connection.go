package database

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/kfsantos/books/models"
)

var DB *gorm.DB

func Setup() {
	db, err := gorm.Open(
		os.Getenv("DB_CONNECTION"),
			"host="+os.Getenv("DB_HOST")+
			" port="+os.Getenv("DB_PORT")+
			" user="+os.Getenv("DB_USERNAME")+
			" dbname="+os.Getenv("DB_DATABASE")+
			" sslmode=disable password="+os.Getenv("DB_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}
	//Para visualizar consultas SQL
	db.LogMode(false)
	//cria a tabela chamada books
	db.AutoMigrate([]models.Book{})
	DB = db
}
