package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"readcommend/models"
	"readcommend/services"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestBookHandler_GetBooks_Success(t *testing.T) {
	expectedBooks := []models.Book{{ID: 1, Title: "Book 1"}}
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	bookService := services.NewBookService(sqlxDB)
	bookHandler := NewBookHandler(bookService)

	req := httptest.NewRequest("GET", "/api/v1/books", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Book 1"))
	bookHandler.GetBooks(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var responseBooks []models.Book
	err := json.Unmarshal(w.Body.Bytes(), &responseBooks)
	assert.NoError(t, err)
	assert.Equal(t, expectedBooks, responseBooks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookHandler_GetBooks_Error(t *testing.T) {
	expectedError := errors.New("database error")
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	bookService := services.NewBookService(sqlxDB)
	bookHandler := NewBookHandler(bookService)

	req := httptest.NewRequest("GET", "/api/v1/books", nil)
	w := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WillReturnError(expectedError)
	bookHandler.GetBooks(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedError.Error(), errorResponse.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestValidateQueryParameters(t *testing.T) {
	testCases := []struct {
		description  string
		queryParams  url.Values
		expectedErr  bool
		expectedData *ParsedQueryParams
	}{
		{
			description: "Valid parameters",
			queryParams: url.Values{
				"authors":   {"1,2,3"},
				"genres":    {"4,5"},
				"min-pages": {"100"},
				"max-pages": {"200"},
				"min-year":  {"1990"},
				"max-year":  {"2020"},
				"limit":     {"10"},
			},
			expectedErr: false,
			expectedData: &ParsedQueryParams{
				Authors:  []int64{1, 2, 3},
				Genres:   []int64{4, 5},
				MinPages: 100,
				MaxPages: 200,
				MinYear:  1990,
				MaxYear:  2020,
				Limit:    10,
			},
		},
		{
			description: "Invalid min-pages",
			queryParams: url.Values{
				"min-pages": {"invalid"},
			},
			expectedErr:  true,
			expectedData: nil,
		},
		{
			description: "Valid parameters with empty values",
			queryParams: url.Values{
				"authors":   {""},
				"genres":    {""},
				"min-pages": {""},
				"max-pages": {""},
				"min-year":  {""},
				"max-year":  {""},
				"limit":     {""},
			},
			expectedErr:  true,
			expectedData: nil,
		},
		{
			description: "Invalid limit (negative value)",
			queryParams: url.Values{
				"limit": {"-10"},
			},
			expectedErr:  true,
			expectedData: nil,
		},
		{
			description: "Valid parameters with limit only",
			queryParams: url.Values{
				"limit": {"5"},
			},
			expectedErr: false,
			expectedData: &ParsedQueryParams{
				Authors:  nil,
				Genres:   nil,
				MinPages: 0,
				MaxPages: 0,
				MinYear:  0,
				MaxYear:  0,
				Limit:    5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			parsedParams, err := ValidateQueryParameters(tc.queryParams, nil)

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, parsedParams)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, parsedParams)
				assert.Equal(t, tc.expectedData, parsedParams)
			}
		})
	}
}
