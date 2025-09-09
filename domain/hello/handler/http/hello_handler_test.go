package http

import (
	"cutbray/first_api/utils/response"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHelloHandler_Hello(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new gin engine
	server := gin.New()

	// Initialize the handler
	NewHelloHandler(server)

	// Test cases
	tests := []struct {
		name            string
		method          string
		path            string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "successful hello request",
			method:          "GET",
			path:            "/",
			expectedStatus:  http.StatusOK,
			expectedMessage: "Hello, World Moncos Lowrider!",
		},
		{
			name:            "successful post body helo request",
			method:          "POST",
			path:            "/post-hello",
			expectedStatus:  http.StatusOK,
			expectedMessage: "Hello, World Moncos Lowrider!",
		},
		{
			name:            "successful post form data helo request",
			method:          "POST",
			path:            "/post-hello-form",
			expectedStatus:  http.StatusOK,
			expectedMessage: "Hello, World Moncos Lowrider!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest(tt.method, tt.path, nil)
			assert.NoError(t, err)

			// Create a response recorder
			w := httptest.NewRecorder()

			// Perform the request
			server.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response body
			var response response.SuccessResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert response structure
			assert.Equal(t, tt.expectedStatus, response.Status)
			// assert.Equal(t, tt.expectedMessage, response.Message)
			assert.NotNil(t, response.Data)
			assert.IsType(t, []interface{}{}, response.Data)
		})
	}
}

func TestHelloHandler_InvalidMethod(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new gin engine
	server := gin.New()

	// Initialize the handler
	NewHelloHandler(server)

	// Test invalid HTTP method
	req, err := http.NewRequest("POST", "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// Should return 404 since POST method is not allowed
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHelloHandler_InvalidPath(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new gin engine
	server := gin.New()

	// Initialize the handler
	NewHelloHandler(server)

	// Test invalid path
	req, err := http.NewRequest("GET", "/invalid", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// Should return 404 for invalid path
	assert.Equal(t, http.StatusNotFound, w.Code)
}
