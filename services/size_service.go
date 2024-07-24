package services

import (
	"fmt"
	"readcommend/models"

	"github.com/jmoiron/sqlx"
)

type SizeService struct {
	DB *sqlx.DB
}

func NewSizeService(db *sqlx.DB) *SizeService {
	return &SizeService{DB: db}
}

func (s *SizeService) GetSizes() ([]models.Size, error) {
	var sizes = []models.Size{}
	err := s.DB.Select(&sizes, "SELECT id, title, min_pages, max_pages FROM size")
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return sizes, nil
}
