package middleware

import (
    "net/http"
    "sync"
    "time"
    "warehouse-api/utils"
    "golang.org/x/time/rate"
)

// Penampung limit per IP
type IPRateLimiter struct {
    ips map[string]*rate.Limiter
    mu  sync.Mutex
    r   rate.Limit
    b   int
}

// Membuat limiter baru
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
    return &IPRateLimiter{
        ips: make(map[string]*rate.Limiter),
        r:   r, // rate (request per second)
        b:   b, // burst size
    }
}

// Mendapatkan limiter untuk IP tertentu, atau membuat baru jika belum ada
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
    i.mu.Lock()
    defer i.mu.Unlock()

    limiter, exists := i.ips[ip]
    if !exists {
        limiter = rate.NewLimiter(i.r, i.b)
        i.ips[ip] = limiter
    }

    return limiter
}

// Variabel global limiter (misal: 5 request per detik, burst 10)
var globalLimiter = NewIPRateLimiter(5, 10)

// Middleware Rate Limit
func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Ambil IP dari RemoteAddr (sederhana)
        // Jika di belakang proxy/load balancer, gunakan X-Forwarded-For
        ip := r.RemoteAddr

        limiter := globalLimiter.GetLimiter(ip)
        if !limiter.Allow() {
            utils.JSONError(w, http.StatusTooManyRequests, "Terlalu banyak permintaan (Rate Limit Exceeded)")
            return
        }

        next.ServeHTTP(w, r)
    })
}

// Cleanup rutin untuk menghapus IP lama (opsional untuk mencegah memory leak)
func init() {
    go func() {
        for {
            time.Sleep(10 * time.Minute)
            globalLimiter.mu.Lock()
            // Reset map setiap 10 menit (sangat sederhana)
            // Implementasi production grade akan mengecek waktu akses terakhir tiap IP
            globalLimiter.ips = make(map[string]*rate.Limiter)
            globalLimiter.mu.Unlock()
        }
    }()
}
