package apifunc

import (
	b "book/storage"

	"github.com/gin-gonic/gin"
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type handler struct {
	storage *b.DBManager
}

// @title           Swagger for book api
// @version         1.0
// @description     This is a book service api.
// @host      		localhost:8000
func NewServer(storage *b.DBManager) *gin.Engine {
	r := gin.Default()

	h := handler{
		storage: storage,
	}

	r.GET("/book/:id", h.GetBook)
	r.POST("/book", h.CreateBook)
	r.PUT("/book/:id", h.UpdateBook)
	r.GET("book/", h.GetAll)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swagFiles.Handler))
	return r
}
