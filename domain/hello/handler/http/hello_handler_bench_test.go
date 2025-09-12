package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkHelloHandler(b *testing.B) {
	server, _ := setupTestServer()
	tokenString, _ := setupTestUser()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", tokenString)

	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
	}
}

func BenchmarkHelloHandler_Parallel(b *testing.B) {
	server, _ := setupTestServer()
	tokenString, _ := setupTestUser()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", tokenString)

	b.ResetTimer()

	// Run the benchmark in parallel
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)
		}
	})
}
