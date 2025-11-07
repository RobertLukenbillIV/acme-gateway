package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	requestsPerSecond := 5
	handler := RateLimit(requestsPerSecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Wrap with RequestID for error handling
	handler = RequestID(handler)

	// Make requests up to the limit
	for i := 0; i < requestsPerSecond; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status %d, got %d", i+1, http.StatusOK, w.Code)
		}
	}

	// Next request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status %d for rate limited request, got %d", http.StatusTooManyRequests, w.Code)
	}

	// Wait for refill
	time.Sleep(1100 * time.Millisecond)

	// Should work again
	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("After refill: expected status %d, got %d", http.StatusOK, w.Code)
	}
}
