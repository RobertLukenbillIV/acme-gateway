package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/RobertLukenbillIV/acme-gateway/internal/errors"
)

// ErrorHandler middleware catches panics and converts them to error responses
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
				requestID := GetRequestID(r.Context())
				errors.WriteError(w, errors.CodeInternalError, "Internal server error", http.StatusInternalServerError, requestID)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
