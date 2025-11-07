package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/RobertLukenbillIV/acme-gateway/internal/errors"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	tokens        int
	maxTokens     int
	refillRate    int
	mu            sync.Mutex
	lastRefillTime time.Time
}

func newRateLimiter(requestsPerSecond int) *RateLimiter {
	return &RateLimiter{
		tokens:         requestsPerSecond,
		maxTokens:      requestsPerSecond,
		refillRate:     requestsPerSecond,
		lastRefillTime: time.Now(),
	}
}

func (rl *RateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(rl.lastRefillTime).Seconds()
	tokensToAdd := int(elapsed * float64(rl.refillRate))

	if tokensToAdd > 0 {
		rl.tokens = min(rl.maxTokens, rl.tokens+tokensToAdd)
		rl.lastRefillTime = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RateLimit middleware enforces rate limiting
func RateLimit(requestsPerSecond int) func(http.Handler) http.Handler {
	limiter := newRateLimiter(requestsPerSecond)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.allow() {
				requestID := GetRequestID(r.Context())
				errors.WriteError(w, errors.CodeRateLimitExceeded, "Rate limit exceeded", http.StatusTooManyRequests, requestID)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
