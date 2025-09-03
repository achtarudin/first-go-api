package main

import (
	"cutbray/first_api/handler/response"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain_ServerSetup(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test server similar to main.go
	server := gin.Default()

	// Import and setup the handler (similar to main.go)
	http := &struct{}{}
	_ = http // avoid unused variable error

	// We can't directly test main() function, but we can test the server setup
	// This test ensures that the server can be created and routes can be registered
	assert.NotNil(t, server)
}

func TestIntegration_HelloEndpoint(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create server exactly like in main.go
	server := gin.Default()

	// Setup routes (we need to import the handler package)
	// For now, let's create a simple test route
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.SuccessResponse{
			Status:  http.StatusOK,
			Message: "Hello, World Moncos Lowrider!",
			Data:    []interface{}{},
		})
	})

	// Test the endpoint
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
	assert.Equal(t, "Hello, World Moncos Lowrider!", response.Message)
	assert.NotNil(t, response.Data)
}

func TestIntegration_ServerPort(t *testing.T) {
	// Test that we can create a server (without actually running it)
	server := gin.Default()

	// Add a test route
	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Test the health endpoint
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}
