package integration

import (
	"log"
	"net/http"
	"readcommend/server"
	"testing"
)

func TestBookAPI(t *testing.T) {
	baseURL := "http://localhost:4001/api/v1"
	client := &http.Client{}

	server, err := server.NewServerInstance("../config.integration.yaml")
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}

	server.Start()
	defer server.Close()

	tests := []struct {
		name               string
		endpoint           string
		method             string
		expectedStatusCode int
	}{
		{
			name:               "Get All Books",
			endpoint:           "/books",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Get Filtered Books",
			endpoint:           "/books?authors=1,2&genres=3,4&min-pages=100&max-pages=300&min-year=2000&max-year=2020&limit=10",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Parameter",
			endpoint:           "/books?invalidparam=abc",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Authors Filter",
			endpoint:           "/books?authors=1,2,3",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Authors Filter Invalid",
			endpoint:           "/books?authors=1,2,invalid",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Genres Filter",
			endpoint:           "/books?genres=3,4,5",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Genres Filter Invalid",
			endpoint:           "/books?genres=3,4,invalid",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Min Pages Filter",
			endpoint:           "/books?min-pages=100",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Min Pages Filter Invalid",
			endpoint:           "/books?min-pages=invalid",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Max Pages Filter",
			endpoint:           "/books?max-pages=500",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Max Pages Filter Invalid",
			endpoint:           "/books?max-pages=invalid",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Min and Max Pages Filter",
			endpoint:           "/books?min-pages=100&max-pages=500",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Min and Max Pages Filter Invalid",
			endpoint:           "/books?min-pages=invalid&max-pages=invalid",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Min Pages larger than Max Pages Filter",
			endpoint:           "/books?min-pages=300&max-pages=200",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Min Year Filter",
			endpoint:           "/books?min-year=2000",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Min Year Filter Invalid Value",
			endpoint:           "/books?min-year=-2000",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Max Year Filter",
			endpoint:           "/books?max-year=2020",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Max Year Filter Invalid",
			endpoint:           "/books?max-year=invalid",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Min and Max Year Filter Invalid",
			endpoint:           "/books?min-year=2010&max-year=2000",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Limit Filter",
			endpoint:           "/books?limit=10",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Limit Filter Invalid",
			endpoint:           "/books?limit=0",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Multiple Filters",
			endpoint:           "/books?authors=1,2&genres=3,4&min-pages=100&max-pages=300&min-year=2000&max-year=2020&limit=10",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Multiple Filters Invalid Values",
			endpoint:           "/books?authors=1,2,invalid&genres=3,4,invalid&min-pages=invalid&max-pages=invalid&min-year=invalid&max-year=invalid&limit=invalid",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := baseURL + test.endpoint
			req, err := http.NewRequest(test.method, url, nil)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error sending request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", test.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}
