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

func TestGenreHandler_GetGenres_Success(t *testing.T) {
	expectedGenres := []models.Genre{{ID: 1, Title: "Genre 1"}}
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	genreService := services.NewGenreService(sqlxDB)
	genreHandler := NewGenreHandler(genreService)

	req := httptest.NewRequest("GET", "/api/v1/genres", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Genre 1"))

	genreHandler.GetGenres(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var responseGenres []models.Genre
	err := json.Unmarshal(w.Body.Bytes(), &responseGenres)
	assert.NoError(t, err)
	assert.Equal(t, expectedGenres, responseGenres)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGenreHandler_GetGenres_Error(t *testing.T) {
	expectedError := errors.New("database error")
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	genreService := services.NewGenreService(sqlxDB)
	genreHandler := NewGenreHandler(genreService)

	req := httptest.NewRequest("GET", "/api/v1/genres", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnError(expectedError)

	genreHandler.GetGenres(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedError.Error(), errorResponse.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
