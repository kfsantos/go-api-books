package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kfsantos/books/database"
	"github.com/kfsantos/books/models"
)

//Estrutura APIEnv que tem objetivo de remover a dependência do objeto DB
type APIEnv struct {
	DB *gorm.DB
}

func (a *APIEnv) GetBooks(c *gin.Context) {
	books, err := database.GetBooks(a.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if len(books) < 1 {
		c.JSON(http.StatusNotFound, "there is no book in db")
		return
	}

	c.JSON(http.StatusOK, books)
}

func (a *APIEnv) GetBook(c *gin.Context) {
	id := c.Params.ByName("id")
	book, exists, err := database.GetBookByID(id, a.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, "there is no book in db")
		return
	}

	c.JSON(http.StatusOK, book)
}

func (a *APIEnv) CreateBook(c *gin.Context) {
	book := models.Book{}
	err := c.BindJSON(&book)
	if err != nil {
		//Verifica se os dados do form estão vazio
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, "Empty fields, no data to save!")
			return
		}		
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := a.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, book)
}

func (a *APIEnv) DeleteBook(c *gin.Context) {
	id := c.Params.ByName("id")
	_, exists, err := database.GetBookByID(id, a.DB)
	if err != nil {
		fmt.Println("Error: not found id------------->", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		fmt.Println("Error: Not found------------->", err)
		c.JSON(http.StatusNotFound, "record not exists")
		return
	} else {
		database.DeleteBook(id, a.DB)
	}
	
	c.JSON(http.StatusOK, "record deleted successfully")
}

func (a *APIEnv) UpdateBook(c *gin.Context) {
	id := c.Params.ByName("id")
	newid, err := strconv.Atoi(id)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, "A valid identifier is required")
		return
	}
	updatedBook := models.Book{}
	err = a.DB.First(&updatedBook, newid).Error
	
	if err != nil {
		c.JSON(http.StatusNotFound, "record not exists")
		return
	}

	err = c.ShouldBindJSON(&updatedBook)
	if err != nil {		
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} 
	err = a.DB.Save(&updatedBook).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to update item!")
		return
	}

	c.JSON(http.StatusOK, updatedBook)
}
