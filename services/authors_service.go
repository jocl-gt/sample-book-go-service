package services

import (
	"fmt"
	"readcommend/models"

	"github.com/jmoiron/sqlx"
)

type AuthorService struct {
	DB *sqlx.DB
}

func NewAuthorService(db *sqlx.DB) *AuthorService {
	return &AuthorService{DB: db}
}

func (s *AuthorService) GetAuthors() ([]models.Author, error) {
	var authors = []models.Author{}
	query := "SELECT id, first_name, last_name FROM author"

	err := s.DB.Select(&authors, query)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}

	return authors, nil
}
