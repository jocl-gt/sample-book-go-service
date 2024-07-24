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

func TestSizeHandler_GetSizes_Success(t *testing.T) {

	minPagesVal := 100
	MaxPagesVal := 200
	expectedSizes := []models.Size{
		{ID: 1, Title: "Size 1", MinPages: nil, MaxPages: nil},
		{ID: 2, Title: "Size 2", MinPages: &minPagesVal, MaxPages: &MaxPagesVal},
	}
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	sizeService := services.NewSizeService(sqlxDB)
	sizeHandler := NewSizeHandler(sizeService)

	req := httptest.NewRequest("GET", "/api/v1/sizes", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "min_pages", "max_pages"}).
		AddRow(1, "Size 1", nil, nil).
		AddRow(2, "Size 2", 100, 200))

	sizeHandler.GetSizes(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var responseSizes []models.Size
	err := json.Unmarshal(w.Body.Bytes(), &responseSizes)
	assert.NoError(t, err)
	assert.Equal(t, expectedSizes, responseSizes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSizeHandler_GetSizes_Error(t *testing.T) {
	expectedError := errors.New("database error")
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	sizeService := services.NewSizeService(sqlxDB)
	sizeHandler := NewSizeHandler(sizeService)

	req := httptest.NewRequest("GET", "/api/v1/sizes", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnError(expectedError)

	sizeHandler.GetSizes(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedError.Error(), errorResponse.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
