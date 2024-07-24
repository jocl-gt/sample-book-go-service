package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"readcommend/services"
	"readcommend/utils"
	"strconv"
)

type ParsedQueryParams struct {
	Authors  []int64
	Genres   []int64
	MinPages int
	MaxPages int
	MinYear  int
	MaxYear  int
	Limit    int
}

type BookHandler struct {
	BookService *services.BookService
}

func NewBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{BookService: bookService}
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	parsedParams, validationErr := ValidateQueryParameters(queryParams, w)
	if validationErr != nil {
		return
	}

	books, err := h.BookService.GetBooks(parsedParams.Authors, parsedParams.Genres, parsedParams.MinPages, parsedParams.MaxPages, parsedParams.MinYear, parsedParams.MaxYear, parsedParams.Limit)
	if err != nil {
		utils.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(books); err != nil {
		utils.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func ValidateQueryParameters(queryParams url.Values, w http.ResponseWriter) (*ParsedQueryParams, error) {
	var parsedParams ParsedQueryParams

	for paramName, paramValue := range queryParams {
		switch paramName {
		case "authors":
			if len(paramValue) > 0 {
				authors, err := utils.ParseQueryParamIntSlice(paramValue[0])
				if err != nil {
					errMsg := fmt.Errorf("invalid value for authors: %s. Please provide valid positive integer author ID(s)", paramValue[0])
					utils.SendErrorResponse(w, errMsg, http.StatusBadRequest)
					return nil, errMsg
				}
				parsedParams.Authors = authors
			}
		case "genres":
			if len(paramValue) > 0 {
				genres, err := utils.ParseQueryParamIntSlice(paramValue[0])
				if err != nil {
					errMsg := fmt.Errorf("invalid value for genres: %s. Please provide valid positive integer genre ID(s)", paramValue[0])
					utils.SendErrorResponse(w, errMsg, http.StatusBadRequest)
					return nil, errMsg
				}
				parsedParams.Genres = genres
			}
		case "min-pages":
			if len(paramValue) > 0 {
				minPages, err := strconv.Atoi(paramValue[0])
				if err != nil || minPages < 0 {
					errMsg := fmt.Errorf("invalid value for min-pages: %s. Please provide a non-negative integer value", paramValue[0])
					utils.SendErrorResponse(w, errMsg, http.StatusBadRequest)
					return nil, errMsg
				}
				parsedParams.MinPages = minPages
			}
		case "max-pages":
			if len(paramValue) > 0 {
				maxPages, err := strconv.Atoi(paramValue[0])
				if err != nil || maxPages < 0 {
					errMsg := fmt.Errorf("invalid value for max-pages: %s. Please provide a non-negative integer value", paramValue[0])
					utils.SendErrorResponse(w, errMsg, http.StatusBadRequest)
					return nil, errMsg
				}
				parsedParams.MaxPages = maxPages
			}
		case "min-year":
			if len(paramValue) > 0 {
				minYear, err := strconv.Atoi(paramValue[0])
				if err != nil || minYear < 0 {
					errMsg := fmt.Errorf("invalid value for min-year: %s. Please provide a non-negative integer value", paramValue[0])
					utils.SendErrorResponse(w, errMsg, http.StatusBadRequest)
					return nil, errMsg
				}
				parsedParams.MinYear = minYear
			}
		case "max-year":
			if len(paramValue) > 0 {
				maxYear, err := strconv.Atoi(paramValue[0])
				if err != nil || maxYear < 0 {
					errMsg := fmt.Errorf("invalid value for max-year: %s. Please provide a non-negative integer value", paramValue[0])
					utils.SendErrorResponse(w, errMsg, http.StatusBadRequest)
					return nil, errMsg
				}
				parsedParams.MaxYear = maxYear
			}
		case "limit":
			if len(paramValue) > 0 {
				limit, err := strconv.Atoi(paramValue[0])
				if err != nil || limit < 1 {
					errMsg := fmt.Errorf("invalid value for limit: %s. Please provide a positive integer value", paramValue[0])
					utils.SendErrorResponse(w, errMsg, http.StatusBadRequest)
					return nil, errMsg
				}
				parsedParams.Limit = limit
			}
		default:
			errMsg := fmt.Errorf("invalid query parameter: %s", paramName)
			utils.SendErrorResponse(w, errMsg, http.StatusBadRequest)
			return nil, errMsg
		}
	}

	if parsedParams.MinPages > 0 && parsedParams.MaxPages > 0 && parsedParams.MinPages > parsedParams.MaxPages {
		utils.SendErrorResponse(w, fmt.Errorf("invalid values for min-pages and max-pages. Min-pages should be less than or equal to max-pages"), http.StatusBadRequest)
		return nil, fmt.Errorf("invalid values for min-pages and max-pages. Min-pages should be less than or equal to max-pages")
	}

	if parsedParams.MinYear > 0 && parsedParams.MaxYear > 0 && parsedParams.MinYear > parsedParams.MaxYear {
		utils.SendErrorResponse(w, fmt.Errorf("invalid values for min-year and max-year. Min-year should be less than or equal to max-year"), http.StatusBadRequest)
		return nil, fmt.Errorf("invalid values for min-year and max-year. Min-year should be less than or equal to max-year")
	}

	return &parsedParams, nil
}
