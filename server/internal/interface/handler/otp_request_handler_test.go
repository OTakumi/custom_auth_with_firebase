package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"

	"custom_auth_api/internal/infrastructure/emailsender"
	"custom_auth_api/internal/infrastructure/persistence"
	"custom_auth_api/internal/interface/handler"
	"custom_auth_api/internal/usecase"
)

const testProjectID = "demo-project"

// setupTestEnvironment initializes Firebase clients and services for testing.
func setupTestEnvironment(t *testing.T) (*firestore.Client, *auth.Client, *handler.OTPRequestHandler, *handler.OTPVerifyHandler, context.Context) {
	t.Helper()

	// Skip if emulator is not running
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("Skipping integration test: FIRESTORE_EMULATOR_HOST is not set.")
	}

	if os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") == "" {
		t.Skip("Skipping integration test: FIREBASE_AUTH_EMULATOR_HOST is not set.")
	}

	ctx := context.Background()

	// Initialize Firebase app
	config := &firebase.Config{ProjectID: testProjectID}

	app, err := firebase.NewApp(ctx, config, option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	// Initialize Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}

	// Initialize Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		t.Fatalf("Failed to create Auth client: %v", err)
	}

	// Initialize services and handlers
	otpRepo := persistence.NewOTPRepository(firestoreClient)
	emailSender := emailsender.NewDummyEmailSender()
	otpService := usecase.NewOTPService(otpRepo, emailSender)
	authService := usecase.NewAuthService(authClient)
	otpRequestHandler := handler.NewOTPRequestHandler(otpService, authService)
	otpVerifyHandler := handler.NewOTPVerifyHandler(otpService, authService)

	return firestoreClient, authClient, otpRequestHandler, otpVerifyHandler, ctx
}

// cleanupOTP deletes the OTP document for the given email.
func cleanupOTP(ctx context.Context, t *testing.T, client *firestore.Client, email string) {
	t.Helper()

	_, err := client.Collection("otps").Doc(email).Delete(ctx)
	if err != nil {
		t.Logf("Failed to clean up OTP data for %s: %v", email, err)
	}
}

// cleanupUser deletes the test user from Firebase Auth.
func cleanupUser(ctx context.Context, t *testing.T, authClient *auth.Client, email string) {
	t.Helper()

	user, err := authClient.GetUserByEmail(ctx, email)
	if err == nil {
		err := authClient.DeleteUser(ctx, user.UID)
		if err != nil {
			t.Logf("Failed to clean up user %s: %v", email, err)
		}
	}
}

// createTestUser creates a test user in Firebase Auth Emulator.
func createTestUser(t *testing.T, authClient *auth.Client, email string) *auth.UserRecord {
	t.Helper()

	ctx := context.Background()
	params := (&auth.UserToCreate{}).
		Email(email).
		Password("password123").
		EmailVerified(true)

	user, err := authClient.CreateUser(ctx, params)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

func TestOTPRequestHandler_RequestOTP_Success(t *testing.T) {
	firestoreClient, authClient, otpRequestHandler, _, ctx := setupTestEnvironment(t)

	defer func() {
		err := firestoreClient.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	email := "request-otp-success@example.com"

	t.Cleanup(func() {
		cleanupOTP(ctx, t, firestoreClient, email)
		cleanupUser(ctx, t, authClient, email)
	})

	// Arrange: Create test user
	createTestUser(t, authClient, email)

	// Create request
	reqBody := map[string]string{"email": email}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/auth/otp", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act: Call handler
	otpRequestHandler.RequestOTP(c)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	message, ok := response["message"].(string)
	if !ok || message == "" {
		t.Error("Expected success message in response")
	}

	// Verify OTP was saved in Firestore
	doc, err := firestoreClient.Collection("otps").Doc(email).Get(ctx)
	if err != nil {
		t.Errorf("Expected OTP to be saved in Firestore, got error: %v", err)
	}

	if doc != nil {
		otp, ok := doc.Data()["otp"].(string)
		if !ok || len(otp) != 6 {
			t.Error("Expected 6-digit OTP to be saved")
		}
	}
}

func TestOTPRequestHandler_RequestOTP_InvalidJSON(t *testing.T) {
	_, _, otpRequestHandler, _, _ := setupTestEnvironment(t) //nolint:dogsled // Only need otpRequestHandler for this test

	// Create request with invalid JSON
	req, _ := http.NewRequest(http.MethodPost, "/auth/otp", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	otpRequestHandler.RequestOTP(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Invalid request body" {
		t.Errorf("Expected 'Invalid request body' error, got %v", response["error"])
	}
}

func TestOTPRequestHandler_RequestOTP_InvalidEmail(t *testing.T) {
	_, _, otpRequestHandler, _, _ := setupTestEnvironment(t) //nolint:dogsled // Only need otpRequestHandler for this test

	// Create request with invalid email
	reqBody := map[string]string{"email": "invalid-email"}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/auth/otp", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	otpRequestHandler.RequestOTP(c)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestOTPRequestHandler_RequestOTP_UserNotFound(t *testing.T) {
	firestoreClient, _, otpRequestHandler, _, _ := setupTestEnvironment(t) //nolint:dogsled // Only need firestoreClient and otpRequestHandler

	defer func() {
		err := firestoreClient.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	email := "nonexistent@example.com"

	// Create request for non-existent user
	reqBody := map[string]string{"email": email}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/auth/otp", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	otpRequestHandler.RequestOTP(c)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Should use generic error message to prevent enumeration
	if response["error"] != "Authentication failed" {
		t.Errorf("Expected 'Authentication failed' error, got %v", response["error"])
	}
}
