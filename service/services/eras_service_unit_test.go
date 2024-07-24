package services

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestErasService_GetEras(t *testing.T) {
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
			expectedRows: sqlmock.NewRows([]string{"id", "title", "min_year", "max_year"}).
				AddRow(1, "Era 1", 1000, 1500).
				AddRow(2, "Era 2", 1500, 1800),
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
			erasService := NewErasService(sqlxDB)

			if tc.expectedRows != nil {
				mock.ExpectQuery("SELECT").WillReturnRows(tc.expectedRows)
			} else {
				mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))
			}
			eras, err := erasService.GetEras()

			assert.NoError(t, mock.ExpectationsWereMet())

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Nil(t, eras)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, eras)
			}
		})
	}
}
