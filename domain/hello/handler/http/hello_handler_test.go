package http

import (
	"cutbray/first_api/domain/auth/entity"
	"cutbray/first_api/pkg/middleware"
	"cutbray/first_api/pkg/response"
	"cutbray/first_api/pkg/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() (*gin.Engine, *helloHandler) {
	gin.SetMode(gin.TestMode)
	server := gin.Default()

	middlewareFunc2 := middleware.JWTAuth()

	handler := NewHelloHandler(server, &middlewareFunc2)
	handler.RegisterRoute()

	return server, handler
}

func setupTestUser() (string, error) {
	user := &entity.User{
		ID:    1,
		Email: "user@example.com",
	}

	return utils.GenerateToken(user)
}
func TestHelloHandler_Hello(t *testing.T) {

	server, handler := setupTestServer()
	assert.NotNil(t, server)
	assert.NotNil(t, handler)

	tokenString, err := setupTestUser()
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Create a request
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", tokenString)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	server.ServeHTTP(w, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response response.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
}
func TestHelloHandler_Hello_Unauthorized(t *testing.T) {

	server, handler := setupTestServer()
	assert.NotNil(t, server)
	assert.NotNil(t, handler)

	// Create a request
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	server.ServeHTTP(w, req)

	// Assert the status code
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Parse response body
	var response response.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, response.Status)
}

func TestHelloHandler_PostBodyHello(t *testing.T) {

	server, handler := setupTestServer()
	assert.NotNil(t, server)
	assert.NotNil(t, handler)

	// Create a request
	req, err := http.NewRequest("POST", "/post-hello", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	server.ServeHTTP(w, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response response.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
}

func TestHelloHandler_PostFormDataHello(t *testing.T) {

	server, handler := setupTestServer()
	assert.NotNil(t, server)
	assert.NotNil(t, handler)

	// Create a request
	req, err := http.NewRequest("POST", "/post-hello-form", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	server.ServeHTTP(w, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response response.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)
}
