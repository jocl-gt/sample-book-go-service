package services

import (
	"fmt"
	"readcommend/models"

	"github.com/jmoiron/sqlx"
)

type ErasService struct {
	DB *sqlx.DB
}

func NewErasService(db *sqlx.DB) *ErasService {
	return &ErasService{DB: db}
}

func (s *ErasService) GetEras() ([]models.Era, error) {
	var eras = []models.Era{}
	err := s.DB.Select(&eras, "SELECT id, title, min_year, max_year FROM era")
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return eras, nil
}
