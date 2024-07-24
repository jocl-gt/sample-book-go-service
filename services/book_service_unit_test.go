package services

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestBookService_GetBooks(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	testCases := []struct {
		name         string
		expectedErr  error
		expectedRows *sqlmock.Rows
	}{
		{
			name:        "Valid request with results",
			expectedErr: nil,
			expectedRows: sqlmock.NewRows([]string{"id", "title", "year_published", "rating", "pages", "genre_id", "author_id"}).
				AddRow(1, "Book 1", 2020, 4.5, 300, 1, 1).
				AddRow(2, "Book 2", 2021, 4.0, 250, 2, 2),
		},
		{
			name:         "Valid request with no results",
			expectedErr:  nil,
			expectedRows: sqlmock.NewRows([]string{}),
		},
		{
			name:         "Database error",
			expectedErr:  errors.New("database error"),
			expectedRows: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bookService := NewBookService(sqlxDB)

			if tc.expectedRows != nil {
				mock.ExpectQuery("SELECT").WillReturnRows(tc.expectedRows)
			} else {
				mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))
			}

			authors := []int64{1, 2}
			genres := []int64{1, 2}
			minPages := 100
			maxPages := 300
			minYear := 2020
			maxYear := 2022
			limit := 10

			books, err := bookService.GetBooks(authors, genres, minPages, maxPages, minYear, maxYear, limit)

			assert.NoError(t, mock.ExpectationsWereMet())

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Nil(t, books)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, books)
			}
		})
	}
}
