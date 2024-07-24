package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"readcommend/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()

	err := errors.New("Sample error message")
	statusCode := http.StatusBadRequest

	SendErrorResponse(w, err, statusCode)
	assert.Equal(t, statusCode, w.Code)

	var response models.ErrorResponse
	errDecode := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, errDecode)

	assert.Equal(t, err.Error(), response.Message)
}

func TestParseQueryParamIntSlice(t *testing.T) {
	testCases := []struct {
		name          string
		paramValue    string
		expectedSlice []int64
		expectedError error
	}{
		{
			name:          "Valid integer values",
			paramValue:    "1,2,3",
			expectedSlice: []int64{1, 2, 3},
			expectedError: nil,
		},
		{
			name:          "Empty string",
			paramValue:    "",
			expectedSlice: []int64(nil),
			expectedError: nil,
		},
		{
			name:          "Invalid integer value",
			paramValue:    "1,2,invalid",
			expectedSlice: nil,
			expectedError: errors.New("strconv.ParseInt: parsing \"invalid\": invalid syntax"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseQueryParamIntSlice(tc.paramValue)

			assert.Equal(t, tc.expectedSlice, result)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
