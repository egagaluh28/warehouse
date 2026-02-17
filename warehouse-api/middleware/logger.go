package middleware

import (
    "log"
    "net/http"
    "time"
)

// Logger Middleware untuk mencatat setiap request yang masuk
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Bungkus ResponseWriter untuk menangkap status code
        rw := &responseWriter{w, http.StatusOK}

        next.ServeHTTP(rw, r)

        duration := time.Since(start)

        // Format Log: [METHOD] URL | Status | Duration
        log.Printf("[%s] %s | %d | %v", r.Method, r.URL.Path, rw.status, duration)
    })
}

// Struct pembungkus untuk menangkap status code response
type responseWriter struct {
    http.ResponseWriter
    status int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}
