package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RobertLukenbillIV/acme-gateway/internal/config"
	"github.com/RobertLukenbillIV/acme-gateway/internal/middleware"
)

func TestRouterTicketsRoute(t *testing.T) {
	// Create a test backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the path was stripped correctly
		if r.URL.Path != "/" && r.URL.Path != "/123" {
			t.Errorf("Expected path / or /123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ticket response"))
	}))
	defer backend.Close()

	cfg := &config.Config{
		TicketsServiceURL: backend.URL,
	}

	router := NewRouter(cfg)
	handler := middleware.RequestID(router)

	// Test /api/tickets route
	req := httptest.NewRequest("GET", "/api/tickets", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Test /api/tickets/123 route
	req = httptest.NewRequest("GET", "/api/tickets/123", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRouterNotFound(t *testing.T) {
	cfg := &config.Config{
		TicketsServiceURL: "http://localhost:9999",
	}

	router := NewRouter(cfg)
	handler := middleware.RequestID(router)

	req := httptest.NewRequest("GET", "/api/unknown", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for unknown route, got %d", http.StatusNotFound, w.Code)
	}
}
