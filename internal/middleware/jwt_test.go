package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTValid(t *testing.T) {
	secret := "test-secret"
	tenantID := "tenant-123"
	roles := []string{"admin", "user"}
	scopes := []string{"read", "write"}

	// Create a valid token
	claims := &Claims{
		TenantID: tenantID,
		Roles:    roles,
		Scopes:   scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	handler := JWT(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify context values
		if GetTenantID(r.Context()) != tenantID {
			t.Errorf("Expected tenant ID %s, got %s", tenantID, GetTenantID(r.Context()))
		}

		// Verify headers
		if r.Header.Get("X-Tenant-ID") != tenantID {
			t.Errorf("Expected X-Tenant-ID header %s, got %s", tenantID, r.Header.Get("X-Tenant-ID"))
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestJWTMissing(t *testing.T) {
	handler := JWT("test-secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without token")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestJWTInvalid(t *testing.T) {
	handler := JWT("test-secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid token")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
