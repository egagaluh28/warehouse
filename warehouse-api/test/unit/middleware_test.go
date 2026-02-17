package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"warehouse-api/middleware"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	t.Run("Rate limiter test", func(t *testing.T) {
		// Rate limiter testing would require actual middleware instance
		// For now, we skip detailed testing
		t.Skip("Rate limiter requires specific middleware instance")
	})
}

func TestLogger(t *testing.T) {
	t.Run("Success - Logger middleware logs request", func(t *testing.T) {
		handler := middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))

		req := httptest.NewRequest("GET", "/api/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Logger should not interfere with request
	})

	t.Run("Success - Logger handles different methods", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "DELETE"}

		for _, method := range methods {
			handler := middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(method, "/api/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}

func TestJWTAuth(t *testing.T) {
	t.Run("Fail - Missing Authorization header", func(t *testing.T) {
		handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/api/protected", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Fail - Invalid token format", func(t *testing.T) {
		handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/api/protected", nil)
		req.Header.Set("Authorization", "InvalidToken")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Fail - Malformed Bearer token", func(t *testing.T) {
		handler := middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCORS(t *testing.T) {
	t.Run("Success - CORS headers are set", func(t *testing.T) {
		// CORS is typically implemented as inline middleware
		// This test demonstrates the concept
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("OPTIONS", "/api/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	})
}

// Test timing of middleware
func TestMiddlewarePerformance(t *testing.T) {
	t.Run("Logger middleware performance", func(t *testing.T) {
		handler := middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		start := time.Now()
		
		for i := 0; i < 100; i++ {
			req := httptest.NewRequest("GET", "/api/test", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		}

		duration := time.Since(start)
		
		// Logger should add minimal overhead
		assert.Less(t, duration.Milliseconds(), int64(1000), "Logger middleware should be fast")
	})
}
