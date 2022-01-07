package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

// GetDB helps you to get a connection
func GetDB() *gorm.DB {
	return DB
}

func GetBooks(db *gorm.DB) ([]models.Book, error) {
	books := []models.Book{}
	err := db.Find(&books).Error

	if err != nil {
		return books, err
	}
	return books, nil
}

func GetBookByID(id string, db *gorm.DB) (models.Book, bool, error) {
	b := models.Book{}

	query := db.Select("books.*")
	query = query.Group("books.id")
	err := query.Where("books.id = ?", id).First(&b).Error

	// err := db.Find(&b).Where(id).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return b, false, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return b, false, nil
	}
	return b, true, nil
}

func DeleteBook(id string, db *gorm.DB) error {
	var b models.Book
	// err := db.Delete(&b).Where(id).Error
	if err := db.Where("id = ?", id).Delete(&b).Error; err != nil {
		// if err != nil {
		return err
	}
	return nil
}

func UpdateBook(db *gorm.DB, b *models.Book) error {
	if err := db.Save(&b).Error; err != nil {
		return err
	}
	return nil
}

func ClearTable() {
	DB.Exec("DELETE FROM books")
	DB.Exec("ALTER SEQUENCE books_id_seq RESTART WITH 1")
}
