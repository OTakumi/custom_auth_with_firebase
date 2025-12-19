package email_test

import (
	"errors"
	"testing"

	email "custom_auth_api/internal/domain/vo"
)

func TestNewEmail(t *testing.T) {
	t.Parallel()
	t.Run("should create a new email for a valid address", func(t *testing.T) {
		t.Parallel()

		validEmail := "test@example.com"

		email, err := email.NewEmail(validEmail)
		if err != nil {
			t.Fatalf("Expected no error for valid email, but got %v", err)
		}

		if email == nil {
			t.Fatal("Expected email to be non-nil for valid email")
		}

		if email.Value != validEmail {
			t.Errorf("Expected email value to be %s, but got %s", validEmail, email.Value)
		}
	})

	testCases := []struct {
		name        string
		email       string
		expectedErr error
	}{
		{
			name:        "should return an error for an email without an @ symbol",
			email:       "test.example.com",
			expectedErr: email.ErrInvalidEmailFormat,
		},
		{
			name:        "should return an error for an email without a domain",
			email:       "test@",
			expectedErr: email.ErrInvalidEmailFormat,
		},
		{
			name:        "should return an error for an email without a user",
			email:       "@example.com",
			expectedErr: email.ErrInvalidEmailFormat,
		},
		{
			name:        "should return an error for an email with invalid characters",
			email:       "test()@example.com",
			expectedErr: email.ErrInvalidEmailFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			email, err := email.NewEmail(tc.email)
			if err == nil {
				t.Errorf("Expected an error for invalid email, but got nil")
			}

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error to be %v, but got %v", tc.expectedErr, err)
			}

			if email != nil {
				t.Errorf("Expected email to be nil for invalid email")
			}
		})
	}
}
