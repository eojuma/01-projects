package middleware

import (
	"net/http"
	"sync"
	"time"
)

// simple in-memory token-bucket per IP rate limiter
func RateLimitMiddleware(next http.Handler, maxRequests int, window time.Duration) http.Handler {
	type bucket struct {
		tokens       int
		lastRefill   time.Time
	}

	buckets := make(map[string]*bucket)
	mu := sync.Mutex{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		mu.Lock()
		b, ok := buckets[ip]
		if !ok {
			b = &bucket{tokens: maxRequests, lastRefill: time.Now()}
			buckets[ip] = b
		}
		// refill
		now := time.Now()
		elapsed := now.Sub(b.lastRefill)
		if elapsed >= window {
			b.tokens = maxRequests
			b.lastRefill = now
		}
		if b.tokens <= 0 {
			mu.Unlock()
			w.Header().Set("Retry-After", "60")
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		b.tokens--
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
