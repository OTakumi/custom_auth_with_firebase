package tests_test

import (
	"context"
	"os"
	"testing"
	"time"

	"custom_auth_api/internal/infrastructure/emailsender"
	"custom_auth_api/internal/infrastructure/persistence"
	"custom_auth_api/internal/usecase"

	"cloud.google.com/go/firestore"
)

// setupIntegrationTest initializes a Firestore client for the emulator.
// It requires the FIRESTORE_EMULATOR_HOST environment variable to be set.
func setupIntegrationTest(t *testing.T) *firestore.Client {
	t.Helper()

	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		t.Skip("Skipping integration test: FIRESTORE_EMULATOR_HOST is not set.")
	}

	// The client will automatically connect to the emulator because the
	// FIRESTORE_EMULATOR_HOST environment variable is set.
	// We need to provide a valid project ID, but it can be any string for the emulator.
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "demo-project")
	if err != nil {
		t.Fatalf("Failed to create Firestore client for emulator: %v", err)
	}

	return client
}

func TestOTPService_Integration_GenerateAndSendOTP(t *testing.T) {
	client := setupIntegrationTest(t)
	t.Cleanup(func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close firestore client: %v", err)
		}
	})

	otpRepo := persistence.NewOTPRepository(client)
	emailSender := emailsender.NewDummyEmailSender()
	otpService := usecase.NewOTPService(otpRepo, emailSender)

	ctx := context.Background()
	testEmail := "integration-test@example.com"

	// Cleanup function to delete the document after the test
	t.Cleanup(func() {
		_, err := client.Collection("otps").Doc(testEmail).Delete(ctx)
		if err != nil {
			t.Logf("Failed to clean up test data: %v", err)
		}
	})

	t.Run("should generate, save, and send OTP", func(t *testing.T) {
		generatedOTP, err := otpService.GenerateAndSendOTP(ctx, testEmail)
		if err != nil {
			t.Fatalf("GenerateAndSendOTP failed: %v", err)
		}

		// Verify the OTP was saved correctly in Firestore
		doc, err := client.Collection("otps").Doc(testEmail).Get(ctx)
		if err != nil {
			t.Fatalf("Failed to get OTP document from Firestore: %v", err)
		}

		data := doc.Data()

		savedOTP, ok := data["otp"].(string)
		if !ok {
			t.Fatal("'otp' field is not a string in Firestore document")
		}

		if savedOTP != generatedOTP {
			t.Errorf("expected saved OTP to be %s, but got %s", generatedOTP, savedOTP)
		}

		expiresAt, ok := data["expiresAt"].(time.Time)
		if !ok {
			t.Fatal("'expiresAt' field is not a time.Time in Firestore document")
		}

		// Check that the expiration time is in the future
		if !expiresAt.After(time.Now()) {
			t.Errorf("expected expiresAt to be in the future, but it was %v", expiresAt)
		}
	})
}
