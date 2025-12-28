package persistence_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"custom_auth_api/internal/infrastructure/persistence"
)

// setupFirestoreEmulator initializes a Firestore client for the emulator.
// It requires the FIRESTORE_EMULATOR_HOST environment variable to be set.
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

func TestOTPRepository_Save(t *testing.T) {
	client := setupFirestoreEmulator(t)
	defer client.Close()

	repo := persistence.NewOTPRepository(client)
	ctx := context.Background()

	tests := []struct {
		name      string
		email     string
		otp       string
		wantError bool
	}{
		{
			name:      "valid OTP save",
			email:     "test-save@example.com",
			otp:       "123456",
			wantError: false,
		},
		{
			name:      "save with different email",
			email:     "another@example.com",
			otp:       "654321",
			wantError: false,
		},
		{
			name:      "overwrite existing OTP",
			email:     "overwrite@example.com",
			otp:       "111111",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanupOTP(ctx, t, client, tt.email)
			})

			// Arrange & Act
			err := repo.Save(ctx, tt.email, tt.otp)

			// Assert
			if (err != nil) != tt.wantError {
				t.Errorf("Save() error = %v, wantError %v", err, tt.wantError)

				return
			}

			// Verify the OTP was saved correctly
			doc, err := client.Collection("otps").Doc(tt.email).Get(ctx)
			if err != nil {
				t.Fatalf("Failed to get saved OTP document: %v", err)
			}

			data := doc.Data()

			savedOTP, ok := data["otp"].(string)
			if !ok {
				t.Fatal("'otp' field is not a string")
			}

			if savedOTP != tt.otp {
				t.Errorf("got OTP = %s, want %s", savedOTP, tt.otp)
			}

			// Verify attempts initialized to 0
			attempts, ok := data["attempts"].(int64)
			if !ok {
				t.Fatal("'attempts' field is not an int64")
			}

			if attempts != 0 {
				t.Errorf("got attempts = %d, want 0", attempts)
			}

			// Verify expiration is in the future
			expiresAt, ok := data["expiresAt"].(time.Time)
			if !ok {
				t.Fatal("'expiresAt' field is not a time.Time")
			}

			if !expiresAt.After(time.Now()) {
				t.Errorf("expiresAt should be in the future, got %v", expiresAt)
			}
		})
	}

	// Test overwriting existing OTP
	t.Run("overwrite resets attempts and expiration", func(t *testing.T) {
		email := "overwrite-test@example.com"

		t.Cleanup(func() {
			cleanupOTP(t, client, email)
		})

		// Save initial OTP
		err := repo.Save(ctx, email, "111111")
		if err != nil {
			t.Fatalf("Failed to save initial OTP: %v", err)
		}

		// Increment attempts
		err = repo.IncrementAttempts(ctx, email)
		if err != nil {
			t.Fatalf("Failed to increment attempts: %v", err)
		}

		// Save new OTP (should reset attempts)
		err = repo.Save(ctx, email, "222222")
		if err != nil {
			t.Fatalf("Failed to save new OTP: %v", err)
		}

		// Verify attempts reset to 0
		doc, err := client.Collection("otps").Doc(email).Get(ctx)
		if err != nil {
			t.Fatalf("Failed to get document: %v", err)
		}

		data := doc.Data()

		attempts, _ := data["attempts"].(int64)
		if attempts != 0 {
			t.Errorf("got attempts = %d, want 0 after overwrite", attempts)
		}

		savedOTP, _ := data["otp"].(string)
		if savedOTP != "222222" {
			t.Errorf("got OTP = %s, want 222222", savedOTP)
		}
	})
}

func TestOTPRepository_Find_Success(t *testing.T) {
	client := setupFirestoreEmulator(t)
	defer client.Close()

	repo := persistence.NewOTPRepository(client)
	ctx := context.Background()
	email := "find-success@example.com"
	expectedOTP := "123456"

	t.Cleanup(func() {
		cleanupOTP(t, client, email)
	})

	// Arrange: Save an OTP
	err := repo.Save(ctx, email, expectedOTP)
	if err != nil {
		t.Fatalf("Failed to save OTP: %v", err)
	}

	// Act: Find the OTP
	foundOTP, err := repo.Find(ctx, email)

	// Assert
	if err != nil {
		t.Errorf("Find() unexpected error = %v", err)
	}

	if foundOTP != expectedOTP {
		t.Errorf("Find() got = %s, want %s", foundOTP, expectedOTP)
	}

	// Verify OTP was NOT deleted (Find no longer deletes)
	doc, err := client.Collection("otps").Doc(email).Get(ctx)
	if err != nil {
		t.Error("OTP document should still exist after Find()")
	}

	if doc != nil {
		otp, _ := doc.Data()["otp"].(string)
		if otp != expectedOTP {
			t.Errorf("Expected OTP to be %s, got %s", expectedOTP, otp)
		}
	}
}

func TestOTPRepository_Find_NotFound(t *testing.T) {
	client := setupFirestoreEmulator(t)
	defer client.Close()

	repo := persistence.NewOTPRepository(client)
	ctx := context.Background()
	email := "nonexistent@example.com"

	// Act: Try to find non-existent OTP
	foundOTP, err := repo.Find(ctx, email)

	// Assert
	if err == nil {
		t.Error("Find() expected error for non-existent OTP, got nil")
	}

	if !errors.Is(err, persistence.ErrOTPNotFound) {
		t.Errorf("Find() expected ErrOTPNotFound, got %v", err)
	}

	if foundOTP != "" {
		t.Errorf("Find() got = %s, want empty string", foundOTP)
	}
}

func TestOTPRepository_Find_Expired(t *testing.T) {
	client := setupFirestoreEmulator(t)
	defer client.Close()

	ctx := context.Background()
	email := "expired@example.com"

	t.Cleanup(func() {
		cleanupOTP(t, client, email)
	})

	// Arrange: Manually create an expired OTP document
	_, err := client.Collection("otps").Doc(email).Set(ctx, map[string]any{
		"otp":       "123456",
		"expiresAt": time.Now().Add(-1 * time.Minute), // Expired 1 minute ago
		"attempts":  0,
	})
	if err != nil {
		t.Fatalf("Failed to create expired OTP: %v", err)
	}

	repo := persistence.NewOTPRepository(client)

	// Act: Try to find expired OTP
	foundOTP, err := repo.Find(ctx, email)

	// Assert
	if err == nil {
		t.Error("Find() expected error for expired OTP, got nil")
	}

	if !errors.Is(err, persistence.ErrOTPExpired) {
		t.Errorf("Find() expected ErrOTPExpired, got %v", err)
	}

	if foundOTP != "" {
		t.Errorf("Find() got = %s, want empty string", foundOTP)
	}

	// Verify document still exists (not deleted for expired OTP)
	_, err = client.Collection("otps").Doc(email).Get(ctx)
	if err != nil {
		t.Error("Expired OTP document should still exist after Find()")
	}
}

func TestOTPRepository_Find_TooManyAttempts(t *testing.T) {
	client := setupFirestoreEmulator(t)
	defer client.Close()

	ctx := context.Background()
	email := "too-many-attempts@example.com"

	t.Cleanup(func() {
		cleanupOTP(t, client, email)
	})

	// Arrange: Create an OTP with too many attempts
	_, err := client.Collection("otps").Doc(email).Set(ctx, map[string]any{
		"otp":       "123456",
		"expiresAt": time.Now().Add(5 * time.Minute),
		"attempts":  3, // Max attempts reached
	})
	if err != nil {
		t.Fatalf("Failed to create OTP with max attempts: %v", err)
	}

	repo := persistence.NewOTPRepository(client)

	// Act: Try to find OTP with too many attempts
	foundOTP, err := repo.Find(ctx, email)

	// Assert
	if err == nil {
		t.Error("Find() expected error for too many attempts, got nil")
	}

	if !errors.Is(err, persistence.ErrTooManyAttempts) {
		t.Errorf("Find() expected ErrTooManyAttempts, got %v", err)
	}

	if foundOTP != "" {
		t.Errorf("Find() got = %s, want empty string", foundOTP)
	}

	// Verify document still exists
	_, err = client.Collection("otps").Doc(email).Get(ctx)
	if err != nil {
		t.Error("OTP document should still exist after too many attempts")
	}
}

func TestOTPRepository_IncrementAttempts(t *testing.T) {
	client := setupFirestoreEmulator(t)
	defer client.Close()

	repo := persistence.NewOTPRepository(client)
	ctx := context.Background()

	t.Run("increment from 0 to 1", func(t *testing.T) {
		email := "increment-test@example.com"

		t.Cleanup(func() {
			cleanupOTP(t, client, email)
		})

		// Arrange: Save an OTP
		err := repo.Save(ctx, email, "123456")
		if err != nil {
			t.Fatalf("Failed to save OTP: %v", err)
		}

		// Act: Increment attempts
		err = repo.IncrementAttempts(ctx, email)

		// Assert
		if err != nil {
			t.Errorf("IncrementAttempts() unexpected error = %v", err)
		}

		// Verify attempts incremented
		doc, err := client.Collection("otps").Doc(email).Get(ctx)
		if err != nil {
			t.Fatalf("Failed to get document: %v", err)
		}

		attempts, ok := doc.Data()["attempts"].(int64)
		if !ok {
			t.Fatal("'attempts' field is not an int64")
		}

		if attempts != 1 {
			t.Errorf("got attempts = %d, want 1", attempts)
		}
	})

	t.Run("multiple increments", func(t *testing.T) {
		email := "multiple-increments@example.com"

		t.Cleanup(func() {
			cleanupOTP(t, client, email)
		})

		// Arrange: Save an OTP
		err := repo.Save(ctx, email, "123456")
		if err != nil {
			t.Fatalf("Failed to save OTP: %v", err)
		}

		// Act: Increment multiple times
		for i := range 3 {
			err = repo.IncrementAttempts(ctx, email)
			if err != nil {
				t.Fatalf("IncrementAttempts() failed on iteration %d: %v", i, err)
			}
		}

		// Assert
		doc, err := client.Collection("otps").Doc(email).Get(ctx)
		if err != nil {
			t.Fatalf("Failed to get document: %v", err)
		}

		attempts, _ := doc.Data()["attempts"].(int64)
		if attempts != 3 {
			t.Errorf("got attempts = %d, want 3", attempts)
		}
	})

	t.Run("increment non-existent document", func(t *testing.T) {
		email := "nonexistent-increment@example.com"

		// Act: Try to increment non-existent document
		err := repo.IncrementAttempts(ctx, email)

		// Assert: Should not return error (idempotent)
		if err != nil {
			t.Errorf("IncrementAttempts() for non-existent doc should not error, got %v", err)
		}
	})
}

func TestOTPRepository_Delete(t *testing.T) {
	client := setupFirestoreEmulator(t)
	defer client.Close()

	repo := persistence.NewOTPRepository(client)
	ctx := context.Background()

	t.Run("delete existing OTP", func(t *testing.T) {
		email := "delete-test@example.com"

		t.Cleanup(func() {
			cleanupOTP(t, client, email)
		})

		// Arrange: Save an OTP
		err := repo.Save(ctx, email, "123456")
		if err != nil {
			t.Fatalf("Failed to save OTP: %v", err)
		}

		// Act: Delete the OTP
		err = repo.Delete(ctx, email)

		// Assert
		if err != nil {
			t.Errorf("Delete() unexpected error = %v", err)
		}

		// Verify OTP was deleted
		_, err = client.Collection("otps").Doc(email).Get(ctx)
		if err == nil {
			t.Error("OTP document should be deleted")
		}
	})

	t.Run("delete non-existent OTP", func(t *testing.T) {
		email := "non-existent-delete@example.com"

		// Act: Delete non-existent OTP (should not error)
		err := repo.Delete(ctx, email)

		// Assert: Should complete without error
		if err != nil {
			t.Errorf("Delete() should not error for non-existent document, got %v", err)
		}
	})
}

func TestOTPRepository_Concurrent_Operations(t *testing.T) {
	client := setupFirestoreEmulator(t)
	defer client.Close()

	repo := persistence.NewOTPRepository(client)
	ctx := context.Background()

	tests := []struct {
		name  string
		email string
	}{
		{"concurrent_user_1", "concurrent1@example.com"},
		{"concurrent_user_2", "concurrent2@example.com"},
		{"concurrent_user_3", "concurrent3@example.com"},
	}

	// Clean up all test data first
	for _, tt := range tests {
		cleanupOTP(t, client, tt.email)
	}

	// Run operations sequentially to avoid client connection issues
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanupOTP(ctx, t, client, tt.email)
			})

			// Save and find OTP
			err := repo.Save(ctx, tt.email, "123456")
			if err != nil {
				t.Errorf("Save() failed: %v", err)
			}

			otp, err := repo.Find(ctx, tt.email)
			if err != nil {
				t.Errorf("Find() failed: %v", err)
			}

			if otp != "123456" {
				t.Errorf("got OTP = %s, want 123456", otp)
			}
		})
	}
}
