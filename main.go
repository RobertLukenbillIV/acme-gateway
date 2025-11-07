package main

import (
	"log"
	"net/http"

	"github.com/RobertLukenbillIV/acme-gateway/internal/config"
	"github.com/RobertLukenbillIV/acme-gateway/internal/middleware"
	"github.com/RobertLukenbillIV/acme-gateway/internal/proxy"
)

func main() {
	cfg := config.Load()

	// Create router with middleware chain
	handler := middleware.Chain(
		proxy.NewRouter(cfg),
		middleware.RequestID,
		middleware.CORS,
		middleware.RateLimit(cfg.RateLimitPerSecond),
		middleware.JWT(cfg.JWTSecret),
		middleware.ErrorHandler,
	)

	addr := ":" + cfg.Port
	log.Printf("Starting acme-gateway on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
