package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port               string
	JWTSecret          string
	RateLimitPerSecond int
	TicketsServiceURL  string
}

// Load reads configuration from environment variables
func Load() *Config {
	rateLimit := 100 // default
	if rl := os.Getenv("RATE_LIMIT_PER_SECOND"); rl != "" {
		if parsed, err := strconv.Atoi(rl); err == nil {
			rateLimit = parsed
		}
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		RateLimitPerSecond: rateLimit,
		TicketsServiceURL:  getEnv("TICKETS_SERVICE_URL", "http://localhost:8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
