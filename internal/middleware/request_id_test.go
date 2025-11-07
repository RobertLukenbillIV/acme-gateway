package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Request ID not found in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Check response header
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header not set")
	}
}

func TestRequestIDWithExisting(t *testing.T) {
	existingID := "test-request-id-123"

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r.Context())
		if requestID != existingID {
			t.Errorf("Expected request ID %s, got %s", existingID, requestID)
		}
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", existingID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") != existingID {
		t.Errorf("Expected X-Request-ID header %s, got %s", existingID, w.Header().Get("X-Request-ID"))
	}
}
