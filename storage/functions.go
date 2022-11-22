package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type Book struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	AuthorName string    `json:"description"`
	Price      float64   `json:"price"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetBooksQueryParams struct {
	AuthorName string
	Title      string
	Page       int32
	Limit      int32
}

type GetBookResult struct {
	Books []*Book
	Count int32
}

type DBManager struct {
	db *sql.DB
}

func NewDBManager(db *sql.DB) *DBManager {
	return &DBManager{db}
}

func (m *DBManager) CreateBook(book *Book) (*Book, error) {
	query := `
		INSERT INTO book (
			title,
			author_name,
			price,
			amount
		) VALUES ($1, $2, $3, $4)
		RETURNING id, title, author_name, price, amount, created_at
	`

	row := m.db.QueryRow(query,
		book.Title,
		book.AuthorName,
		book.Price,
		book.Amount,
	)

	var result Book

	err := row.Scan(
		&result.ID,
		&result.Title,
		&result.AuthorName,
		&result.Price,
		&result.Amount,
		&result.CreatedAt,
	)
	if err != nil {
		fmt.Println("failed to scan")
		return nil, err
	}

	return &result, nil
}

func (m *DBManager) GetBook(id int) (*Book, error) {
	query := `
		SELECT 
			id,
			title,
			author_name,
			price,
			amount
		FROM book
		WHERE id=$1
	`

	row := m.db.QueryRow(query, id)

	var result Book

	err := row.Scan(
		&result.ID,
		&result.Title,
		&result.AuthorName,
		&result.Price,
		&result.Amount,
	)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *DBManager) UpdateBook(book *Book) (*Book, error) {
	query := `
		UPDATE book SET 
			title=$1,
			author_name=$2,
			price=$3,
			amount=$4
		WHERE id=$5
		RETURNING id, title, author_name, price, amount, created_at
	`
	row := m.db.QueryRow(query,
		book.Title,
		book.AuthorName,
		book.Price,
		book.Amount,
		book.ID,
	)

	var result Book

	err := row.Scan(
		&result.ID,
		&result.Title,
		&result.AuthorName,
		&result.Price,
		&result.Amount,
		&result.CreatedAt,
	)

	if err != nil {
		fmt.Println("salom")
		return nil, err
	}
	return &result, err
}

func (m *DBManager) DeleteBook(id int) error {
	query := `
		DELETE FROM book WHERE id=$1
	`

	_, err := m.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBManager) GetAllBooks(params *GetBooksQueryParams) (*GetBookResult, error) {

	result := GetBookResult{
		Books: make([]*Book, 0),
	}

	offset := (params.Page - 1) * params.Limit
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := "WHERE TRUE"

	if params.AuthorName != "" {
		filter += " AND author_name ilike '%" + params.AuthorName + "%' "
	}

	if params.Title != "" {
		filter += " AND title ilike '%" + params.Title + "%' "
	}

	query := `
		SELECT 
			id,
			title,
			author_name,
			price,
			amount,
			created_at
		FROM book
	` + filter + `
	ORDER BY created_at desc
	` + limit

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.AuthorName,
			&book.Price,
			&book.Amount,
			&book.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result.Books = append(result.Books, &book)
	}
	return &result, nil
}
