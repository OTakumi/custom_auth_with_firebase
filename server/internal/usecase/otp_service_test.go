package usecase_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"custom_auth_api/internal/infrastructure/emailsender"
	"custom_auth_api/internal/infrastructure/persistence"
	"custom_auth_api/internal/usecase"
)

// setupFirestoreEmulator initializes a Firestore client for the emulator.
func setupFirestoreEmulator(t *testing.T) *firestore.Client {
	t.Helper()

	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("Skipping integration test: FIRESTORE_EMULATOR_HOST is not set.")
	}

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "demo-project")
	if err != nil {
		t.Fatalf("Failed to create Firestore client for emulator: %v", err)
	}

	return client
}

// cleanupOTP deletes the OTP document for the given email.
func cleanupOTP(ctx context.Context, t *testing.T, client *firestore.Client, email string) {
	t.Helper()

	_, err := client.Collection("otps").Doc(email).Delete(ctx)
	if err != nil {
		t.Logf("Failed to clean up test data for %s: %v", email, err)
	}
}

func TestOTPService_GenerateAndSendOTP(t *testing.T) {
	client := setupFirestoreEmulator(t)

	defer func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	repo := persistence.NewOTPRepository(client)
	sender := emailsender.NewDummyEmailSender()
	service := usecase.NewOTPService(repo, sender)

	ctx := context.Background()

	tests := []struct {
		name      string
		email     string
		wantError bool
	}{
		{
			name:      "successful OTP generation",
			email:     "test@example.com",
			wantError: false,
		},
		{
			name:      "generate for different email",
			email:     "another@example.com",
			wantError: false,
		},
		{
			name:      "regenerate OTP for same email",
			email:     "regenerate@example.com",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanupOTP(ctx, t, client, tt.email)
			})

			// Act: Generate and send OTP
			otp, err := service.GenerateAndSendOTP(ctx, tt.email)

			// Assert
			if (err != nil) != tt.wantError {
				t.Errorf("GenerateAndSendOTP() error = %v, wantError %v", err, tt.wantError)

				return
			}

			if !tt.wantError {
				// Verify OTP is 6 digits
				if len(otp) != 6 {
					t.Errorf("Expected OTP length 6, got %d", len(otp))
				}

				// Verify OTP is numeric
				for _, c := range otp {
					if c < '0' || c > '9' {
						t.Errorf("OTP contains non-numeric character: %c", c)
					}
				}

				// Verify OTP was saved in Firestore
				doc, err := client.Collection("otps").Doc(tt.email).Get(ctx)
				if err != nil {
					t.Fatalf("Failed to get OTP document from Firestore: %v", err)
				}

				data := doc.Data()

				savedOTP, ok := data["otp"].(string)
				if !ok {
					t.Fatal("'otp' field is not a string")
				}

				if savedOTP != otp {
					t.Errorf("Saved OTP = %s, returned OTP = %s", savedOTP, otp)
				}

				// Verify attempts initialized to 0
				attempts, ok := data["attempts"].(int64)
				if !ok {
					t.Fatal("'attempts' field is not an int64")
				}

				if attempts != 0 {
					t.Errorf("Expected attempts = 0, got %d", attempts)
				}

				// Verify expiration is in the future
				expiresAt, ok := data["expiresAt"].(time.Time)
				if !ok {
					t.Fatal("'expiresAt' field is not a time.Time")
				}

				if !expiresAt.After(time.Now()) {
					t.Errorf("expiresAt should be in the future, got %v", expiresAt)
				}
			}
		})
	}
}

func TestOTPService_VerifyOTP_Success(t *testing.T) {
	client := setupFirestoreEmulator(t)

	defer func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	repo := persistence.NewOTPRepository(client)
	sender := emailsender.NewDummyEmailSender()
	service := usecase.NewOTPService(repo, sender)

	ctx := context.Background()
	email := "verify-success@example.com"

	t.Cleanup(func() {
		cleanupOTP(ctx, t, client, email)
	})

	// Arrange: Generate an OTP
	generatedOTP, err := service.GenerateAndSendOTP(ctx, email)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Act: Verify the OTP
	valid, err := service.VerifyOTP(ctx, email, generatedOTP)

	// Assert
	if err != nil {
		t.Errorf("VerifyOTP() unexpected error = %v", err)
	}

	if !valid {
		t.Error("VerifyOTP() expected valid = true, got false")
	}

	// Verify OTP was deleted after successful verification (one-time use)
	_, err = client.Collection("otps").Doc(email).Get(ctx)
	if err == nil {
		t.Error("OTP should be deleted after successful verification")
	}
}

func TestOTPService_VerifyOTP_InvalidOTP(t *testing.T) {
	client := setupFirestoreEmulator(t)

	defer func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	repo := persistence.NewOTPRepository(client)
	sender := emailsender.NewDummyEmailSender()
	service := usecase.NewOTPService(repo, sender)

	ctx := context.Background()
	email := "verify-invalid@example.com"

	t.Cleanup(func() {
		cleanupOTP(ctx, t, client, email)
	})

	// Arrange: Generate an OTP
	generatedOTP, err := service.GenerateAndSendOTP(ctx, email)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Act: Verify with wrong OTP
	wrongOTP := "999999"
	if wrongOTP == generatedOTP {
		wrongOTP = "888888" // Ensure it's different
	}

	valid, err := service.VerifyOTP(ctx, email, wrongOTP)

	// Assert
	if err == nil {
		t.Error("VerifyOTP() expected error for invalid OTP, got nil")
	}

	if !errors.Is(err, usecase.ErrInvalidOTP) {
		t.Errorf("VerifyOTP() expected ErrInvalidOTP, got %v", err)
	}

	if valid {
		t.Error("VerifyOTP() expected valid = false, got true")
	}

	// Verify attempts was incremented
	doc, err := client.Collection("otps").Doc(email).Get(ctx)
	if err != nil {
		t.Fatalf("Failed to get OTP document: %v", err)
	}

	attempts, ok := doc.Data()["attempts"].(int64)
	if !ok {
		t.Fatal("'attempts' field is not an int64")
	}

	if attempts != 1 {
		t.Errorf("Expected attempts = 1 after failed verification, got %d", attempts)
	}
}

func TestOTPService_VerifyOTP_ExpiredOTP(t *testing.T) {
	client := setupFirestoreEmulator(t)

	defer func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	repo := persistence.NewOTPRepository(client)
	sender := emailsender.NewDummyEmailSender()
	service := usecase.NewOTPService(repo, sender)

	ctx := context.Background()
	email := "verify-expired@example.com"

	t.Cleanup(func() {
		cleanupOTP(ctx, t, client, email)
	})

	// Arrange: Manually create an expired OTP
	expiredOTP := "123456"

	_, err := client.Collection("otps").Doc(email).Set(ctx, map[string]any{
		"otp":       expiredOTP,
		"expiresAt": time.Now().Add(-1 * time.Minute), // Expired 1 minute ago
		"attempts":  0,
	})
	if err != nil {
		t.Fatalf("Failed to create expired OTP: %v", err)
	}

	// Act: Try to verify expired OTP
	valid, err := service.VerifyOTP(ctx, email, expiredOTP)

	// Assert
	if err == nil {
		t.Error("VerifyOTP() expected error for expired OTP, got nil")
	}

	// The error should be from the repository layer (ErrOTPExpired)
	if !errors.Is(err, persistence.ErrOTPExpired) {
		t.Errorf("VerifyOTP() expected error containing ErrOTPExpired, got %v", err)
	}

	if valid {
		t.Error("VerifyOTP() expected valid = false for expired OTP, got true")
	}
}

func TestOTPService_VerifyOTP_TooManyAttempts(t *testing.T) {
	client := setupFirestoreEmulator(t)

	defer func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	repo := persistence.NewOTPRepository(client)
	sender := emailsender.NewDummyEmailSender()
	service := usecase.NewOTPService(repo, sender)

	ctx := context.Background()
	email := "verify-too-many@example.com"

	t.Cleanup(func() {
		cleanupOTP(ctx, t, client, email)
	})

	// Arrange: Create an OTP with max attempts
	validOTP := "123456"

	_, err := client.Collection("otps").Doc(email).Set(ctx, map[string]any{
		"otp":       validOTP,
		"expiresAt": time.Now().Add(5 * time.Minute),
		"attempts":  3, // Max attempts reached
	})
	if err != nil {
		t.Fatalf("Failed to create OTP with max attempts: %v", err)
	}

	// Act: Try to verify OTP with too many attempts
	valid, err := service.VerifyOTP(ctx, email, validOTP)

	// Assert
	if err == nil {
		t.Error("VerifyOTP() expected error for too many attempts, got nil")
	}

	if !errors.Is(err, persistence.ErrTooManyAttempts) {
		t.Errorf("VerifyOTP() expected error containing ErrTooManyAttempts, got %v", err)
	}

	if valid {
		t.Error("VerifyOTP() expected valid = false for too many attempts, got true")
	}
}

func TestOTPService_VerifyOTP_NotFound(t *testing.T) {
	client := setupFirestoreEmulator(t)

	defer func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	repo := persistence.NewOTPRepository(client)
	sender := emailsender.NewDummyEmailSender()
	service := usecase.NewOTPService(repo, sender)

	ctx := context.Background()
	email := "nonexistent@example.com"

	// Act: Try to verify non-existent OTP
	valid, err := service.VerifyOTP(ctx, email, "123456")

	// Assert
	if err == nil {
		t.Error("VerifyOTP() expected error for non-existent OTP, got nil")
	}

	if !errors.Is(err, persistence.ErrOTPNotFound) {
		t.Errorf("VerifyOTP() expected error containing ErrOTPNotFound, got %v", err)
	}

	if valid {
		t.Error("VerifyOTP() expected valid = false for non-existent OTP, got true")
	}
}

func TestOTPService_VerifyOTP_MultipleFailedAttempts(t *testing.T) {
	client := setupFirestoreEmulator(t)

	defer func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	repo := persistence.NewOTPRepository(client)
	sender := emailsender.NewDummyEmailSender()
	service := usecase.NewOTPService(repo, sender)

	ctx := context.Background()
	email := "multiple-attempts@example.com"

	t.Cleanup(func() {
		cleanupOTP(ctx, t, client, email)
	})

	// Arrange: Generate an OTP
	correctOTP, err := service.GenerateAndSendOTP(ctx, email)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Act: Try to verify with wrong OTP multiple times
	wrongOTP := "999999"
	if wrongOTP == correctOTP {
		wrongOTP = "888888"
	}

	// First failed attempt
	valid, err := service.VerifyOTP(ctx, email, wrongOTP)
	if err == nil || valid {
		t.Error("First attempt should fail")
	}

	// Verify attempts = 1
	doc, _ := client.Collection("otps").Doc(email).Get(ctx)

	attempts1, _ := doc.Data()["attempts"].(int64)
	if attempts1 != 1 {
		t.Errorf("Expected attempts = 1, got %d", attempts1)
	}

	// Second failed attempt
	valid, err = service.VerifyOTP(ctx, email, wrongOTP)
	if err == nil || valid {
		t.Error("Second attempt should fail")
	}

	// Verify attempts = 2
	doc, _ = client.Collection("otps").Doc(email).Get(ctx)

	attempts2, _ := doc.Data()["attempts"].(int64)
	if attempts2 != 2 {
		t.Errorf("Expected attempts = 2, got %d", attempts2)
	}

	// Third failed attempt
	valid, err = service.VerifyOTP(ctx, email, wrongOTP)
	if err == nil || valid {
		t.Error("Third attempt should fail")
	}

	// Verify attempts = 3
	doc, _ = client.Collection("otps").Doc(email).Get(ctx)

	attempts3, _ := doc.Data()["attempts"].(int64)
	if attempts3 != 3 {
		t.Errorf("Expected attempts = 3, got %d", attempts3)
	}

	// Fourth attempt should be blocked (too many attempts)
	valid, err = service.VerifyOTP(ctx, email, correctOTP) // Even with correct OTP
	if err == nil {
		t.Error("Fourth attempt should be blocked")
	}

	if !errors.Is(err, persistence.ErrTooManyAttempts) {
		t.Errorf("Expected ErrTooManyAttempts, got %v", err)
	}

	if valid {
		t.Error("VerifyOTP() should return false when attempts exceeded")
	}
}

func TestOTPService_VerifyOTP_SuccessAfterFailedAttempt(t *testing.T) {
	client := setupFirestoreEmulator(t)

	defer func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close Firestore client: %v", err)
		}
	}()

	repo := persistence.NewOTPRepository(client)
	sender := emailsender.NewDummyEmailSender()
	service := usecase.NewOTPService(repo, sender)

	ctx := context.Background()
	email := "success-after-fail@example.com"

	t.Cleanup(func() {
		cleanupOTP(ctx, t, client, email)
	})

	// Arrange: Generate an OTP
	correctOTP, err := service.GenerateAndSendOTP(ctx, email)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Act: First attempt with wrong OTP
	wrongOTP := "999999"
	if wrongOTP == correctOTP {
		wrongOTP = "888888"
	}

	valid, err := service.VerifyOTP(ctx, email, wrongOTP)
	if err == nil || valid {
		t.Error("First attempt with wrong OTP should fail")
	}

	// Second attempt with correct OTP
	valid, err = service.VerifyOTP(ctx, email, correctOTP)

	// Assert
	if err != nil {
		t.Errorf("VerifyOTP() with correct OTP should succeed, got error: %v", err)
	}

	if !valid {
		t.Error("VerifyOTP() with correct OTP should return true")
	}

	// Verify OTP was deleted
	_, err = client.Collection("otps").Doc(email).Get(ctx)
	if err == nil {
		t.Error("OTP should be deleted after successful verification")
	}
}
