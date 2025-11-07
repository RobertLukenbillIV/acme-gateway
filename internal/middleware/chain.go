package middleware

import "net/http"

// Chain applies middleware in the order they are provided
// The first middleware in the chain is the outermost (executed first)
func Chain(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	// Apply middleware in reverse order so the first one wraps all others
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}
