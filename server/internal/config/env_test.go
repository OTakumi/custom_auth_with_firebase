package config

import (
	"os"
	"testing"
)

func TestLoadEnv_Success(t *testing.T) {
	t.Run("loads with default values in development mode", func(t *testing.T) {
		// Arrange - clear environment
		clearEnv(t)

		// Act
		env, err := LoadEnv()

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if env.Port != "8000" {
			t.Errorf("expected port 8000, got %s", env.Port)
		}
		if env.Environment != "development" {
			t.Errorf("expected environment development, got %s", env.Environment)
		}
		if env.RateLimitRequestsPerMinute != 5 {
			t.Errorf("expected 5 requests per minute, got %d", env.RateLimitRequestsPerMinute)
		}
		if env.RateLimitCleanupIntervalMinutes != 10 {
			t.Errorf("expected 10 minute cleanup interval, got %d", env.RateLimitCleanupIntervalMinutes)
		}
	})

	t.Run("loads custom values from environment variables", func(t *testing.T) {
		// Arrange
		clearEnv(t)
		t.Setenv("PORT", "9000")
		t.Setenv("ENV", "production")
		t.Setenv("ALLOWED_ORIGINS", "https://example.com,https://app.example.com")
		t.Setenv("RATE_LIMIT_REQUESTS_PER_MINUTE", "10")
		t.Setenv("RATE_LIMIT_CLEANUP_INTERVAL_MINUTES", "20")

		// Act
		env, err := LoadEnv()

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if env.Port != "9000" {
			t.Errorf("expected port 9000, got %s", env.Port)
		}
		if env.Environment != "production" {
			t.Errorf("expected environment production, got %s", env.Environment)
		}
		if len(env.AllowedOrigins) != 2 {
			t.Errorf("expected 2 allowed origins, got %d", len(env.AllowedOrigins))
		}
		if env.AllowedOrigins[0] != "https://example.com" {
			t.Errorf("expected first origin https://example.com, got %s", env.AllowedOrigins[0])
		}
		if env.RateLimitRequestsPerMinute != 10 {
			t.Errorf("expected 10 requests per minute, got %d", env.RateLimitRequestsPerMinute)
		}
		if env.RateLimitCleanupIntervalMinutes != 20 {
			t.Errorf("expected 20 minute cleanup interval, got %d", env.RateLimitCleanupIntervalMinutes)
		}
	})
}

func TestLoadEnv_ProductionValidation(t *testing.T) {
	t.Run("returns error when ALLOWED_ORIGINS is missing in production", func(t *testing.T) {
		// Arrange
		clearEnv(t)
		t.Setenv("ENV", "production")

		// Act
		env, err := LoadEnv()

		// Assert
		if err == nil {
			t.Error("expected error for missing ALLOWED_ORIGINS in production")
		}
		if env != nil {
			t.Error("expected nil env when error occurs")
		}
		if err.Error() != "ALLOWED_ORIGINS environment variable is required in production" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("loads successfully with ALLOWED_ORIGINS in production", func(t *testing.T) {
		// Arrange
		clearEnv(t)
		t.Setenv("ENV", "production")
		t.Setenv("ALLOWED_ORIGINS", "https://example.com")

		// Act
		env, err := LoadEnv()

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(env.AllowedOrigins) != 1 {
			t.Errorf("expected 1 allowed origin, got %d", len(env.AllowedOrigins))
		}
	})
}

func TestLoadEnv_InvalidRateLimitValues(t *testing.T) {
	t.Run("returns error when RATE_LIMIT_REQUESTS_PER_MINUTE is not an integer", func(t *testing.T) {
		// Arrange
		clearEnv(t)
		t.Setenv("RATE_LIMIT_REQUESTS_PER_MINUTE", "invalid")

		// Act
		env, err := LoadEnv()

		// Assert
		if err == nil {
			t.Error("expected error for invalid RATE_LIMIT_REQUESTS_PER_MINUTE")
		}
		if env != nil {
			t.Error("expected nil env when error occurs")
		}
	})

	t.Run("returns error when RATE_LIMIT_CLEANUP_INTERVAL_MINUTES is not an integer", func(t *testing.T) {
		// Arrange
		clearEnv(t)
		t.Setenv("RATE_LIMIT_CLEANUP_INTERVAL_MINUTES", "not-a-number")

		// Act
		env, err := LoadEnv()

		// Assert
		if err == nil {
			t.Error("expected error for invalid RATE_LIMIT_CLEANUP_INTERVAL_MINUTES")
		}
		if env != nil {
			t.Error("expected nil env when error occurs")
		}
	})
}

func TestEnv_IsProduction(t *testing.T) {
	t.Parallel()

	t.Run("returns true when environment is production", func(t *testing.T) {
		t.Parallel()
		env := &Env{Environment: "production"}
		if !env.IsProduction() {
			t.Error("expected IsProduction to return true")
		}
	})

	t.Run("returns false when environment is not production", func(t *testing.T) {
		t.Parallel()
		env := &Env{Environment: "development"}
		if env.IsProduction() {
			t.Error("expected IsProduction to return false")
		}
	})
}

func TestEnv_IsDevelopment(t *testing.T) {
	t.Parallel()

	t.Run("returns true when environment is development", func(t *testing.T) {
		t.Parallel()
		env := &Env{Environment: "development"}
		if !env.IsDevelopment() {
			t.Error("expected IsDevelopment to return true")
		}
	})

	t.Run("returns false when environment is not development", func(t *testing.T) {
		t.Parallel()
		env := &Env{Environment: "production"}
		if env.IsDevelopment() {
			t.Error("expected IsDevelopment to return false")
		}
	})
}

// clearEnv clears all environment variables used by the config package.
// This ensures tests are isolated and don't interfere with each other.
func clearEnv(t *testing.T) {
	t.Helper()
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("ENV")
	_ = os.Unsetenv("ALLOWED_ORIGINS")
	_ = os.Unsetenv("RATE_LIMIT_REQUESTS_PER_MINUTE")
	_ = os.Unsetenv("RATE_LIMIT_CLEANUP_INTERVAL_MINUTES")
}
