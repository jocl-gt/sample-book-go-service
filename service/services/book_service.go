package services

import (
	"fmt"
	"readcommend/models"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type BookService struct {
	DB *sqlx.DB
}

func NewBookService(db *sqlx.DB) *BookService {
	return &BookService{DB: db}
}

func (s *BookService) GetBooks(authors []int64, genres []int64, minPages int, maxPages int, minYear int, maxYear int, limit int) ([]models.Book, error) {
	query := "SELECT id, title, year_published, rating, pages, genre_id, author_id FROM book"
	var conditions []string
	var args []interface{}

	if len(authors) > 0 {
		authorArray := pq.Int64Array(authors)
		queryPart, authorArgs, err := sqlx.In("author_id = ANY(?)", authorArray)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, queryPart)
		args = append(args, authorArgs...)
	}

	if len(genres) > 0 {
		genreArray := pq.Int64Array(genres)
		queryPart, genreArgs, err := sqlx.In("genre_id = ANY(?)", genreArray)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, queryPart)
		args = append(args, genreArgs...)
	}

	if minPages > 0 {
		conditions = append(conditions, "pages >= ?")
		args = append(args, minPages)
	}

	if maxPages > 0 {
		conditions = append(conditions, "pages <= ?")
		args = append(args, maxPages)
	}

	if minYear > 0 {
		conditions = append(conditions, "year_published >= ?")
		args = append(args, minYear)
	}

	if maxYear > 0 {
		conditions = append(conditions, "year_published <= ?")
		args = append(args, maxYear)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY rating DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	var books = []models.Book{}
	query = s.DB.Rebind(query)
	err := s.DB.Select(&books, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}

	return books, nil
}
