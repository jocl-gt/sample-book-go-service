package utils

import (
	"encoding/json"
	"net/http"
	"readcommend/models"
	"strconv"
	"strings"
)

func SendErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	if w != nil {
		data := models.ErrorResponse{Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if errEncode := json.NewEncoder(w).Encode(data); errEncode != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ParseQueryParamIntSlice(paramValue string) ([]int64, error) {
	var result []int64
	if paramValue != "" {
		values := strings.Split(paramValue, ",")
		for _, value := range values {
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			result = append(result, intValue)
		}
	}
	return result, nil
}
