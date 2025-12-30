package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Configuration errors.
var (
	ErrAllowedOriginsRequired = errors.New("ALLOWED_ORIGINS environment variable is required in production")
	ErrInvalidIntegerValue    = errors.New("environment variable must be a valid integer")
)

// Default configuration values.
const (
	defaultPort                            = "8000"
	defaultEnvironment                     = "development"
	defaultRateLimitRequestsPerMinute      = 5
	defaultRateLimitCleanupIntervalMinutes = 10
)

// Env holds all environment-based configuration values.
// This struct is populated from environment variables and validated on load.
type Env struct {
	// Server configuration
	Port string

	// Environment mode (development/production)
	Environment string

	// CORS configuration
	AllowedOrigins []string

	// Rate limiting configuration
	RateLimitRequestsPerMinute      int
	RateLimitCleanupIntervalMinutes int
}

// LoadEnv loads and validates all environment variables.
// Returns an error if required environment variables are missing or invalid.
func LoadEnv() (*Env, error) {
	env := &Env{
		Port:                            getEnvOrDefault("PORT", defaultPort),
		Environment:                     getEnvOrDefault("ENV", defaultEnvironment),
		AllowedOrigins:                  nil, // Will be set below for production
		RateLimitRequestsPerMinute:      0,   // Will be set below
		RateLimitCleanupIntervalMinutes: 0,   // Will be set below
	}

	// Validate and load CORS origins
	if env.Environment == "production" {
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			return nil, ErrAllowedOriginsRequired
		}
		env.AllowedOrigins = strings.Split(allowedOrigins, ",")
	}
	// In development, AllowedOrigins will be empty and handled by CORS middleware

	// Load rate limiting configuration with defaults
	requestsPerMinute, err := getEnvAsInt("RATE_LIMIT_REQUESTS_PER_MINUTE", defaultRateLimitRequestsPerMinute)
	if err != nil {
		return nil, err
	}
	env.RateLimitRequestsPerMinute = requestsPerMinute

	cleanupInterval, err := getEnvAsInt("RATE_LIMIT_CLEANUP_INTERVAL_MINUTES", defaultRateLimitCleanupIntervalMinutes)
	if err != nil {
		return nil, err
	}
	env.RateLimitCleanupIntervalMinutes = cleanupInterval

	return env, nil
}

// IsProduction returns true if the environment is set to production.
func (e *Env) IsProduction() bool {
	return e.Environment == "production"
}

// IsDevelopment returns true if the environment is set to development.
func (e *Env) IsDevelopment() bool {
	return e.Environment == "development"
}

// getEnvOrDefault retrieves an environment variable or returns a default value.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value.
// Returns an error if the value is not a valid integer.
func getEnvAsInt(key string, defaultValue int) (int, error) {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidIntegerValue, key)
	}

	return value, nil
}
