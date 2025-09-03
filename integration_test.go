package main

import (
	"cutbray/first_api/handler/response"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handlerHttp "cutbray/first_api/handler/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFullIntegration_RealHandler(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create server exactly like in main.go
	server := gin.Default()

	// Use the real handler from your code
	handlerHttp.NewHelloHandler(server)

	// Test the actual endpoint
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response response.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.Status)
	// assert.Equal(t, "Hello, World Moncos Lowrider!", response.Message)
	assert.NotNil(t, response.Data)
	assert.IsType(t, []interface{}{}, response.Data)
}

func TestFullIntegration_ContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := gin.Default()
	handlerHttp.NewHelloHandler(server)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// Check content type is JSON
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestFullIntegration_ResponseHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := gin.Default()
	handlerHttp.NewHelloHandler(server)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// Basic response validation
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, len(w.Body.String()) > 0, "Response body should not be empty")
}
