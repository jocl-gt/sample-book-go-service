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

func TestErasHandler_GetEras_Success(t *testing.T) {
	expectedEras := []models.Era{{ID: 1, Title: "Era 1"}}
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	erasService := services.NewErasService(sqlxDB)
	erasHandler := NewErasHandler(erasService)

	req := httptest.NewRequest("GET", "/api/v1/eras", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Era 1"))

	erasHandler.GetEras(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var responseEras []models.Era
	err := json.Unmarshal(w.Body.Bytes(), &responseEras)
	assert.NoError(t, err)
	assert.Equal(t, expectedEras, responseEras)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestErasHandler_GetEras_Error(t *testing.T) {
	expectedError := errors.New("database error")
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	erasService := services.NewErasService(sqlxDB)
	erasHandler := NewErasHandler(erasService)

	req := httptest.NewRequest("GET", "/api/v1/eras", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnError(expectedError)
	erasHandler.GetEras(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedError.Error(), errorResponse.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
