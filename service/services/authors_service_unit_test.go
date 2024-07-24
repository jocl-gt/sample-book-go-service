package services

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuthorService_GetAuthors(t *testing.T) {
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
			expectedRows: sqlmock.NewRows([]string{"id", "first_name", "last_name"}).
				AddRow(1, "John", "Doe").
				AddRow(2, "Jane", "Smith"),
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
			authorService := NewAuthorService(sqlxDB)

			if tc.expectedRows != nil {
				mock.ExpectQuery("SELECT").WillReturnRows(tc.expectedRows)
			} else {
				mock.ExpectQuery("SELECT").WillReturnError(errors.New("db error"))
			}
			authors, err := authorService.GetAuthors()

			assert.NoError(t, mock.ExpectationsWereMet())

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Nil(t, authors)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, authors)
			}
		})
	}
}
