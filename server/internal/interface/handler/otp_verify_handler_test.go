package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"custom_auth_api/internal/infrastructure/emailsender"
	"custom_auth_api/internal/infrastructure/persistence"
	"custom_auth_api/internal/usecase"
)

func TestOTPVerifyHandler_VerifyOTP_Success(t *testing.T) {
	firestoreClient, authClient, _, otpVerifyHandler, ctx := setupTestEnvironment(t)
	defer firestoreClient.Close()

	email := "verify-otp-success@example.com"

	t.Cleanup(func() {
		cleanupOTP(t, firestoreClient, email)
		cleanupUser(t, authClient, email)
	})

	// Arrange: Create test user
	createTestUser(t, authClient, email, "password123")

	// Generate OTP
	otpRepo := persistence.NewOTPRepository(firestoreClient)
	emailSender := emailsender.NewDummyEmailSender()
	otpService := usecase.NewOTPService(otpRepo, emailSender)

	generatedOTP, err := otpService.GenerateAndSendOTP(ctx, email)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Create request
	reqBody := map[string]string{
		"email": email,
		"otp":   generatedOTP,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	otpVerifyHandler.VerifyOTP(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	token, ok := response["token"].(string)
	if !ok || token == "" {
		t.Error("Expected custom token in response")
	}
}

func TestOTPVerifyHandler_VerifyOTP_InvalidJSON(t *testing.T) {
	_, _, _, otpVerifyHandler, _ := setupTestEnvironment(t)

	// Create request with invalid JSON
	req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	otpVerifyHandler.VerifyOTP(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestOTPVerifyHandler_VerifyOTP_InvalidEmail(t *testing.T) {
	_, _, _, otpVerifyHandler, _ := setupTestEnvironment(t)

	// Create request with invalid email
	reqBody := map[string]string{
		"email": "invalid-email",
		"otp":   "123456",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	otpVerifyHandler.VerifyOTP(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestOTPVerifyHandler_VerifyOTP_InvalidOTP(t *testing.T) {
	firestoreClient, authClient, _, otpVerifyHandler, ctx := setupTestEnvironment(t)
	defer firestoreClient.Close()

	email := "verify-otp-invalid@example.com"

	t.Cleanup(func() {
		cleanupOTP(t, firestoreClient, email)
		cleanupUser(t, authClient, email)
	})

	// Arrange: Create test user and generate OTP
	createTestUser(t, authClient, email, "password123")

	otpRepo := persistence.NewOTPRepository(firestoreClient)
	emailSender := emailsender.NewDummyEmailSender()
	otpService := usecase.NewOTPService(otpRepo, emailSender)

	_, err := otpService.GenerateAndSendOTP(ctx, email)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Create request with wrong OTP
	reqBody := map[string]string{
		"email": email,
		"otp":   "999999",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	otpVerifyHandler.VerifyOTP(c)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "Invalid or expired OTP" {
		t.Errorf("Expected 'Invalid or expired OTP' error, got %v", response["error"])
	}
}

func TestOTPVerifyHandler_VerifyOTP_ExpiredOTP(t *testing.T) {
	firestoreClient, authClient, _, otpVerifyHandler, ctx := setupTestEnvironment(t)
	defer firestoreClient.Close()

	email := "verify-otp-expired@example.com"

	t.Cleanup(func() {
		cleanupOTP(t, firestoreClient, email)
		cleanupUser(t, authClient, email)
	})

	// Arrange: Create test user
	createTestUser(t, authClient, email, "password123")

	// Manually create an expired OTP
	_, err := firestoreClient.Collection("otps").Doc(email).Set(ctx, map[string]any{
		"otp":       "123456",
		"expiresAt": ctx.Value("now"), // Already expired
		"attempts":  0,
	})
	if err != nil {
		t.Fatalf("Failed to create expired OTP: %v", err)
	}

	// Create request
	reqBody := map[string]string{
		"email": email,
		"otp":   "123456",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	otpVerifyHandler.VerifyOTP(c)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestOTPVerifyHandler_VerifyOTP_UserNotFound(t *testing.T) {
	firestoreClient, _, _, otpVerifyHandler, ctx := setupTestEnvironment(t)
	defer firestoreClient.Close()

	email := "verify-user-not-found@example.com"

	t.Cleanup(func() {
		cleanupOTP(t, firestoreClient, email)
	})

	// Arrange: Create OTP but no user
	otpRepo := persistence.NewOTPRepository(firestoreClient)

	err := otpRepo.Save(ctx, email, "123456")
	if err != nil {
		t.Fatalf("Failed to save OTP: %v", err)
	}

	// Create request
	reqBody := map[string]string{
		"email": email,
		"otp":   "123456",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/auth/verify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	otpVerifyHandler.VerifyOTP(c)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)

	// Should use generic error message
	if response["error"] != "Authentication failed" {
		t.Errorf("Expected 'Authentication failed' error, got %v", response["error"])
	}
}
