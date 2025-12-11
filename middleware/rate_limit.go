package middleware

import (
	"net/http"
	"sync"
	"time"

	"reveil-api/config"
	"reveil-api/utils"

	"golang.org/x/time/rate"
)

// RateLimiter manages rate limits for users/IPs
type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
// r: requests per second
// b: burst size
func NewRateLimiter(r float64, b int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     rate.Limit(r),
		burst:    b,
	}
}

// getLimiter returns the limiter for a given key (IP or Token)
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[key] = limiter
	}

	return limiter
}

// CleanOldVisitors removes old entries to prevent memory leak
// Run this in a goroutine
func (rl *RateLimiter) CleanupLoop(interval time.Duration) {
	for {
		time.Sleep(interval)
		rl.mu.Lock()
		// Logic to remove old visitors would require tracking last seen time.
		// For MVP we skip complex cleanup or just wipe map periodically if it gets too big.
		if len(rl.visitors) > 10000 {
			rl.visitors = make(map[string]*rate.Limiter)
		}
		rl.mu.Unlock()
	}
}

// LimitMiddleware applies the rate limit
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prefer User ID if authenticated, else IP
		key := r.RemoteAddr
		userID := r.Context().Value("user_id")
		if userID != nil {
			key = userID.(string)
		}

		limiter := rl.getLimiter(key)
		if !limiter.Allow() {
			utils.ErrorResponseWithCode(w, http.StatusTooManyRequests, "Rate limit exceeded", config.ErrorInternal) // Use better error code
			return
		}

		next.ServeHTTP(w, r)
	})
}
