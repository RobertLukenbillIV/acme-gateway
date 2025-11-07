package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RobertLukenbillIV/acme-gateway/internal/errors"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims from acme-auth-service
type Claims struct {
	TenantID string   `json:"tenant_id"`
	Roles    []string `json:"roles"`
	Scopes   []string `json:"scopes"`
	jwt.RegisteredClaims
}

// JWT middleware validates JWT tokens and extracts claims
func JWT(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				requestID := GetRequestID(r.Context())
				errors.WriteError(w, errors.CodeUnauthorized, "Missing authorization header", http.StatusUnauthorized, requestID)
				return
			}

			// Check for Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				requestID := GetRequestID(r.Context())
				errors.WriteError(w, errors.CodeUnauthorized, "Invalid authorization header format", http.StatusUnauthorized, requestID)
				return
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				requestID := GetRequestID(r.Context())
				errors.WriteError(w, errors.CodeInvalidToken, "Invalid or expired token", http.StatusUnauthorized, requestID)
				return
			}

			// Extract claims
			if claims, ok := token.Claims.(*Claims); ok {
				ctx := r.Context()
				ctx = context.WithValue(ctx, TenantIDKey, claims.TenantID)
				ctx = context.WithValue(ctx, RolesKey, claims.Roles)
				ctx = context.WithValue(ctx, ScopesKey, claims.Scopes)

				// Add headers for downstream services
				r.Header.Set("X-Tenant-ID", claims.TenantID)
				if len(claims.Roles) > 0 {
					r.Header.Set("X-Roles", strings.Join(claims.Roles, ","))
				}
				if len(claims.Scopes) > 0 {
					r.Header.Set("X-Scopes", strings.Join(claims.Scopes, ","))
				}

				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				requestID := GetRequestID(r.Context())
				errors.WriteError(w, errors.CodeInvalidToken, "Invalid token claims", http.StatusUnauthorized, requestID)
			}
		})
	}
}
