package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"custom_auth_api/internal/config"
	"custom_auth_api/internal/interface/handler"
	"custom_auth_api/internal/interface/router"
	"custom_auth_api/internal/usecase"
)

func TestNewRouter_HealthCheck(t *testing.T) {
	t.Parallel()

	// Arrange
	env := &config.Env{
		Port:                            "8000",
		Environment:                     "development",
		RateLimitRequestsPerMinute:      5,
		RateLimitCleanupIntervalMinutes: 10,
	}

	// Create mock handlers (nil services for health check test)
	handlers := &router.Handlers{
		OTPRequest: handler.NewOTPRequestHandler(nil, nil),
		OTPVerify:  handler.NewOTPVerifyHandler(nil, nil),
	}

	r := router.NewRouter(env, handlers)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	expectedBody := `{"message":"OK"}`
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, w.Body.String())
	}
}

func TestNewRouter_RoutesRegistered(t *testing.T) {
	t.Parallel()

	// Arrange
	env := &config.Env{
		Environment:                     "development",
		RateLimitRequestsPerMinute:      5,
		RateLimitCleanupIntervalMinutes: 10,
	}

	// Create mock auth service
	mockAuthService := usecase.NewAuthService(nil)
	handlers := &router.Handlers{
		OTPRequest: handler.NewOTPRequestHandler(nil, mockAuthService),
		OTPVerify:  handler.NewOTPVerifyHandler(nil, mockAuthService),
	}

	r := router.NewRouter(env, handlers)

	testCases := []struct {
		name       string
		method     string
		path       string
		shouldFind bool
	}{
		{
			name:       "health check endpoint exists",
			method:     http.MethodGet,
			path:       "/health",
			shouldFind: true,
		},
		{
			name:       "OTP request endpoint exists",
			method:     http.MethodPost,
			path:       "/auth/otp",
			shouldFind: true,
		},
		{
			name:       "OTP verify endpoint exists",
			method:     http.MethodPost,
			path:       "/auth/verify",
			shouldFind: true,
		},
		{
			name:       "non-existent endpoint returns 404",
			method:     http.MethodGet,
			path:       "/non-existent",
			shouldFind: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Assert
			if tc.shouldFind {
				// For existing endpoints, we expect either 200 or 400 (not 404)
				// 400 might occur if request body is missing, but route is found
				if w.Code == http.StatusNotFound {
					t.Errorf("expected route to exist, but got 404")
				}
			} else {
				if w.Code != http.StatusNotFound {
					t.Errorf("expected 404 for non-existent route, got %d", w.Code)
				}
			}
		})
	}
}
