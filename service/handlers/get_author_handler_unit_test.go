package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"readcommend/models"
	"readcommend/services"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuthorHandler_GetAuthors_Success(t *testing.T) {
	expectedAuthors := []models.Author{{ID: 1, FirstName: "John", LastName: "Wick"}}
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authorService := services.NewAuthorService(sqlxDB)
	authorHandler := NewAuthorHandler(authorService)

	req := httptest.NewRequest("GET", "/api/v1/authors", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name"}).AddRow(1, "John", "Wick"))

	authorHandler.GetAuthors(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var responseAuthors []models.Author
	err := json.Unmarshal(w.Body.Bytes(), &responseAuthors)
	assert.NoError(t, err)
	assert.Equal(t, expectedAuthors, responseAuthors)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthorHandler_GetAuthors_Error(t *testing.T) {
	expectedError := errors.New("database error")
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authorService := services.NewAuthorService(sqlxDB)
	authorHandler := NewAuthorHandler(authorService)

	req := httptest.NewRequest("GET", "/api/v1/authors", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnError(expectedError)
	authorHandler.GetAuthors(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var errorResponse models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedError.Error(), errorResponse.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
