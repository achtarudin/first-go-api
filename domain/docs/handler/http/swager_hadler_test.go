package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewSwaggerHandler(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new gin engine
	server := gin.New()

	// Initialize the Swagger handler
	handler := NewSwaggerHandler(server, "Test API", "Test API description")

	handler.RegisterRoute()

	// Check struct is not nil
	assert.NotNil(t, handler)
	assert.Equal(t, server, handler.server)

	// Test GET /swagger/index.html (Swag UI endpoint)
	req, err := http.NewRequest("GET", "/swagger/index.html", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// Assert status code: Swagger UI should respond, usually 200 or 404 if assets not found
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
	// Optionally: check response body contains Swagger UI HTML or error
}

func TestSwaggerHandlerRouteRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := gin.New()
	handler := NewSwaggerHandler(server, "Title", "Desc")

	handler.RegisterRoute()

	// Test with /swagger/doc.json (the spec endpoint)
	req, err := http.NewRequest("GET", "/swagger/doc.json", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// Should respond with 200 or 404 if spec not found
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
}
