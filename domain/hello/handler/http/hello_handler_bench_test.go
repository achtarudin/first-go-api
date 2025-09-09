package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkHelloHandler(b *testing.B) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new gin engine
	server := gin.New()

	// Initialize the handler
	NewHelloHandler(server)

	// Create a request that we'll reuse
	req, _ := http.NewRequest("GET", "/", nil)

	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
	}
}

func BenchmarkHelloHandler_Parallel(b *testing.B) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new gin engine
	server := gin.New()

	// Initialize the handler
	NewHelloHandler(server)

	// Create a request that we'll reuse
	req, _ := http.NewRequest("GET", "/", nil)

	b.ResetTimer()

	// Run the benchmark in parallel
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)
		}
	})
}
