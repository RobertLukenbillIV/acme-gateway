package errors

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a unified error response matching acme-contracts schemas
type ErrorResponse struct {
	Error   ErrorDetail `json:"error"`
	TraceID string      `json:"trace_id,omitempty"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Standard error codes
const (
	CodeUnauthorized      = "UNAUTHORIZED"
	CodeForbidden         = "FORBIDDEN"
	CodeNotFound          = "NOT_FOUND"
	CodeBadRequest        = "BAD_REQUEST"
	CodeInternalError     = "INTERNAL_ERROR"
	CodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	CodeInvalidToken      = "INVALID_TOKEN"
)

// WriteError writes a standardized error response
func WriteError(w http.ResponseWriter, code string, message string, status int, traceID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Status:  status,
		},
		TraceID: traceID,
	}

	json.NewEncoder(w).Encode(resp)
}
