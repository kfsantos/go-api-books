package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kfsantos/books/database"
	"github.com/kfsantos/books/models"
	"github.com/stretchr/testify/assert"
)

func Test_GetBooks_Ok(t *testing.T) {
	database.Setup()
	db := database.GetDB()
	//Inserindo dados para posterior comparação
	insertTestBook(db)

	req, w := setGetBooksRouter(db)
	defer db.Close()

	a := assert.New(t)
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	actual := models.Book{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	expected := models.Book{}
	a.Equal(expected, actual)
	database.ClearTable()
}

func setGetBooksRouter(db *gorm.DB) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()
	api := &APIEnv{DB: db}
	r.GET("/", api.GetBooks)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func Test_GetBooks_EmptyResult(t *testing.T) {
	database.Setup()
	db := database.GetDB()
	req, w := setGetBooksRouter(db)
	defer db.Close()

	a := assert.New(t)
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request status code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	actual := models.Book{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	expected := models.Book{}
	a.Equal(expected, actual)
	database.ClearTable()
}

func Test_GetBook_InvalidId(t *testing.T) {
	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	req, w := setGetBookRouter(db, "/aa")
	defer db.Close()

	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusInternalServerError, w.Code, "HTTP request status code error")
}

func Test_GetBook_IdNotFound(t *testing.T) {
	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	req, w := setGetBookRouter(db, "/-2")
	defer db.Close()

	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request status code error")
}

func Test_GetBook_OK(t *testing.T) {
	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	book, err := insertTestBook(db)
	if err != nil {
		a.Error(err)
	}

	req, w := setGetBookRouter(db, "/1")
	defer db.Close()

	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	actual := models.Book{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	actual.Model = gorm.Model{}
	expected := book
	expected.Model = gorm.Model{}
	a.Equal(expected, actual)
	database.ClearTable()
}

func setGetBookRouter(db *gorm.DB, url string) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()
	api := &APIEnv{DB: db}
	r.GET("/:id", api.GetBook)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func insertTestBook(db *gorm.DB) (models.Book, error) {
	b := models.Book{
		Author:    "test",
		Name:      "test",
		PageCount: 10,
	}

	if err := db.Create(&b).Error; err != nil {
		return b, err
	}

	return b, nil
}

func Test_CreateBook_OK(t *testing.T) {
	a := assert.New(t)
	database.Setup()
	db := database.GetDB()
	book := models.Book{
		Author:    "test",
		Name:      "test",
		PageCount: 10,
	}

	reqBody, err := json.Marshal(book)
	if err != nil {
		a.Error(err)
	}

	req, w, err := setCreateBookRouter(db, bytes.NewBuffer(reqBody))
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	actual := models.Book{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	actual.Model = gorm.Model{}
	expected := book
	a.Equal(expected, actual)
	database.ClearTable()
}

func setCreateBookRouter(db *gorm.DB,
	body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder, error) {
	r := gin.New()
	api := &APIEnv{DB: db}
	r.POST("/", api.CreateBook)
	req, err := http.NewRequest(http.MethodPost, "/", body)
	if err != nil {
		return req, httptest.NewRecorder(), err
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w, nil
}

func Test_CreateBook_ErrorBind(t *testing.T) {
	type BookTest struct {
		// gorm.Model
		ID        int
		Author    int
		Name      string
		PageCount int
	}
	bookTest := BookTest{
		Author:    3,
		Name:      "test",
		PageCount: 10,
	}

	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	reqBody, err := json.Marshal(bookTest)
	if err != nil {
		a.Error(err)
	}

	req, w, err := setCreateBookRouter(db, bytes.NewBuffer(reqBody))
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusBadRequest, w.Code, "HTTP request status code error")
}

func Test_CreateBook_ErrorCreate(t *testing.T) {
	type BookTestCreate struct {
		// gorm.Model
		ID        int
		Author    string
		Name      string
		PageCount int
	}
	bookTestCreate := BookTestCreate{
		ID:        1,
		Author:    "test",
		Name:      "test",
		PageCount: 10,
	}

	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	insertTestBook(db)

	reqBody, err := json.Marshal(bookTestCreate)
	if err != nil {
		a.Error(err)
	}

	req, w, err := setCreateBookRouter(db, bytes.NewBuffer(reqBody))
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusInternalServerError, w.Code, "HTTP request status code error")
}

func Test_UpdateBook_Ok(t *testing.T) {
	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	bookTestUpdate := models.Book{
		Author:    "test",
		Name:      "test",
		PageCount: 10,
	}
	reqBody1, err1 := json.Marshal(bookTestUpdate)
	if err1 != nil {
		a.Error(err1)
	}

	setCreateBookRouter(db, bytes.NewBuffer(reqBody1))
	defer db.Close()

	bookTestUpdate.ID = 1

	reqBody, err := json.Marshal(bookTestUpdate)
	if err != nil {
		a.Error(err)
	}

	req, w, err := setUpdateBookRouter(db, "/1", bytes.NewBuffer(reqBody))
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodPut, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	actual := models.Book{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	expected := bookTestUpdate
	a.NotEqual(expected, actual)
	database.ClearTable()
}

func setUpdateBookRouter(db *gorm.DB, url string,
	body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder, error) {
	r := gin.New()
	api := &APIEnv{DB: db}

	r.PUT("/:id", api.UpdateBook)
	req, err := http.NewRequest(http.MethodPut, url, body)

	if err != nil {
		return req, httptest.NewRecorder(), err
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return req, w, nil
}

func TestA_UpdateBook_IdNotFound(t *testing.T) {
	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	bookTestUpdate := models.Book{
		Author:    "dsdsds",
		Name:      "dsds",
		PageCount: 10,
	}

	bookTestUpdate.ID = 2

	reqBody, err := json.Marshal(bookTestUpdate)
	if err != nil {
		a.Error(err)
	}

	req, w, err := setUpdateBookRouter(db, "/2", bytes.NewBuffer(reqBody))
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodPut, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request status code error")

	database.ClearTable()
}

func Test_UpdateBook_InvalidId(t *testing.T) {

	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	bookTestUpdate := models.Book{
		Author:    "dsdsds",
		Name:      "dsds",
		PageCount: 10,
	}

	bookTestUpdate.ID = 2

	reqBody, err := json.Marshal(bookTestUpdate)
	if err != nil {
		a.Error(err)
	}

	req, w, err := setUpdateBookRouter(db, "/aa", bytes.NewBuffer(reqBody))
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodPut, req.Method, "HTTP request method error")
	a.Equal(http.StatusInternalServerError, w.Code, "HTTP request status code error")

	database.ClearTable()
}

func Test_UpdateBook_CantShouldBindJSON(t *testing.T) {

	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	bookTestUpdate := models.Book{
		Author:    "test",
		Name:      "test",
		PageCount: 10,
	}
	reqBody1, err1 := json.Marshal(bookTestUpdate)
	if err1 != nil {
		a.Error(err1)
	}

	setCreateBookRouter(db, bytes.NewBuffer(reqBody1))
	defer db.Close()

	reqBody, err := json.Marshal("Teste")
	if err != nil {
		a.Error(err)
	}

	req, w, err := setUpdateBookRouter(db, "/1", bytes.NewBuffer(reqBody))
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodPut, req.Method, "HTTP request method error")
	a.Equal(http.StatusInternalServerError, w.Code, "HTTP request status code error")
	database.ClearTable()
}

func Test_UpdateBook_UpdateError(t *testing.T) {

}

func Test_DeleteBook_Ok(t *testing.T) {
	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	book, err := insertTestBook(db)
	if err != nil {
		a.Error(err)
	}

	req, w, err := setDeleteBookRouter(db, "/1")
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodDelete, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	actual := models.Book{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	actual.Model = gorm.Model{}
	expected := book
	expected.Model = gorm.Model{}
	a.NotEqual(expected, actual)
	database.ClearTable()
}

func setDeleteBookRouter(db *gorm.DB, url string) (*http.Request, *httptest.ResponseRecorder, error) {
	r := gin.New()
	api := &APIEnv{DB: db}
	r.DELETE("/:id", api.DeleteBook)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return req, httptest.NewRecorder(), err
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w, nil
}

func Test_DeleteBook_InvalidId(t *testing.T) {
	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	req, w, err := setDeleteBookRouter(db, "/aa")
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodDelete, req.Method, "HTTP request method error")
	a.Equal(http.StatusInternalServerError, w.Code, "HTTP request status code error")
}

func Test_DeleteBook_IdNotFound(t *testing.T) {
	a := assert.New(t)
	database.Setup()
	db := database.GetDB()

	req, w, err := setDeleteBookRouter(db, "/-1")
	if err != nil {
		a.Error(err)
	}
	defer db.Close()

	a.Equal(http.MethodDelete, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request status code error")
}
