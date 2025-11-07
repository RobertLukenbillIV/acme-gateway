package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/RobertLukenbillIV/acme-gateway/internal/config"
	"github.com/RobertLukenbillIV/acme-gateway/internal/errors"
	"github.com/RobertLukenbillIV/acme-gateway/internal/middleware"
)

// Router handles routing requests to backend services
type Router struct {
	routes map[string]*httputil.ReverseProxy
}

// NewRouter creates a new router with configured service mappings
func NewRouter(cfg *config.Config) *Router {
	router := &Router{
		routes: make(map[string]*httputil.ReverseProxy),
	}

	// Configure route mappings
	if cfg.TicketsServiceURL != "" {
		ticketsURL, err := url.Parse(cfg.TicketsServiceURL)
		if err == nil {
			router.routes["/api/tickets"] = httputil.NewSingleHostReverseProxy(ticketsURL)
		}
	}

	return router
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Find matching route
	for prefix, proxy := range rt.routes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			// Strip the prefix and forward to backend service
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}

			// Forward the request
			proxy.ServeHTTP(w, r)
			return
		}
	}

	// No route found
	requestID := middleware.GetRequestID(r.Context())
	errors.WriteError(w, errors.CodeNotFound, "Route not found", http.StatusNotFound, requestID)
}
