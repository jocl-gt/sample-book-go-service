package services

import (
	"fmt"
	"readcommend/models"

	"github.com/jmoiron/sqlx"
)

type GenreService struct {
	DB *sqlx.DB
}

func NewGenreService(db *sqlx.DB) *GenreService {
	return &GenreService{DB: db}
}

func (s *GenreService) GetGenres() ([]models.Genre, error) {
	var genres = []models.Genre{}
	query := "SELECT id, title FROM genre"

	err := s.DB.Select(&genres, query)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}

	return genres, nil
}
