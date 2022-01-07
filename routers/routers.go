package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kfsantos/books/database"
	"github.com/kfsantos/books/handlers"
)

func Setup() *gin.Engine {
	r := gin.Default()
	api := &handlers.APIEnv{
		DB: database.GetDB(),
	}

	r.GET("", api.GetBooks)
	r.GET("/:id", api.GetBook)
	r.POST("", api.CreateBook)
	r.PUT("/:id", api.UpdateBook)
	r.DELETE("/:id", api.DeleteBook)

	return r
}
