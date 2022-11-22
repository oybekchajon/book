package apifunc

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"book/storage"
)

// @Router /book/{id} [get]
// @Summary Get book by id
// @Description Get book by id
// @Tags book
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} storage.Book
// @Failure 500 {object} ResponseError
func (h *handler) GetBook(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	blog, err := h.storage.GetBook(int(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, blog)
}

// @Router /book [post]
// @Summary Create a book
// @Description Create a book
// @Tags book
// @Accept json
// @Produce json
// @Param book body BookRequest true "Book"
// @Success 200 {object} storage.Book
// @Failure 500 {object} ResponseError
func (h *handler) CreateBook(ctx *gin.Context) {
	var b BookRequest
	if err := ctx.ShouldBindJSON(&b); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	created, err := h.storage.CreateBook(&storage.Book{
		AuthorName: b.Author,
		Title:      b.Title,
		Amount:     b.Amount,
		Price:      b.Price,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, created)
}

// @Router /book/{id} [put]
// @Summary Update a book
// @Description Update a book
// @Tags book
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param book body BookRequest true "Book"
// @Success 200 {object} BookRequest
// @Failure 500 {object} ResponseError
func (h *handler) UpdateBook(ctx *gin.Context) {
	var b BookRequest
	if err := ctx.ShouldBindJSON(&b); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	updated, err := h.storage.UpdateBook(&storage.Book{
		ID:         id,
		Title:      b.Title,
		AuthorName: b.Author,
		Price:      b.Price,
		Amount:     b.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

// @Router /book [get]
// @Summary Get all books
// @Description Get all books
// @Tags book
// @Accept json
// @Produce json
// @Param limit query int true "Limit"
// @Param page query int true "Page"
// @Param author_name query string false "AuthorName"
// @Param title query string false "Title"
// @Success 200 {object} BookRequest
// @Failure 500 {object} ResponseError
func (h *handler) GetAll(ctx *gin.Context) {
	queryParams, err := validateGetBooksQuery(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message" : err.Error(),
		})
		return
	}

	resp, err := h.storage.GetAllBooks(queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message" : err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func validateGetBooksQuery(ctx *gin.Context) (*storage.GetBooksQueryParams, error) {
	var (
		limit int64 = 5
		page  int64 = 1
		err   error
	)

	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &storage.GetBooksQueryParams{
		Limit: int32(limit),
		Page: int32(page),
		AuthorName: ctx.Query("author_name"),
		Title: ctx.Query("title"),
	},nil
}
