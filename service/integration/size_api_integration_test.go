package integration

import (
	"log"
	"net/http"
	"readcommend/server"
	"testing"
)

func TestSizeAPI(t *testing.T) {
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
			name:               "Get All Sizes",
			endpoint:           "/sizes",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid parameter",
			endpoint:           "/sizes?param=1",
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
