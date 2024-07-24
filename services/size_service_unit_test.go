package services

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSizeService_GetSizes(t *testing.T) {
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
			expectedRows: sqlmock.NewRows([]string{"id", "title", "min_pages", "max_pages"}).
				AddRow(1, "Size 1", 100, 200).
				AddRow(2, "Size 2", 200, 300),
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
			sizeService := NewSizeService(sqlxDB)

			if tc.expectedRows != nil {
				mock.ExpectQuery("SELECT").WillReturnRows(tc.expectedRows)
			} else {
				mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))
			}
			sizes, err := sizeService.GetSizes()

			assert.NoError(t, mock.ExpectationsWereMet())

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Nil(t, sizes)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sizes)
			}
		})
	}
}
