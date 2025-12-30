package otp_test

import (
	"regexp"
	"testing"

	"custom_auth_api/internal/domain/vo/otp"
)

func TestNewOTP(t *testing.T) {
	t.Parallel()

	t.Run("should create a new OTP with valid properties", func(t *testing.T) {
		t.Parallel()

		otp, err := otp.NewOTP()
		if err != nil {
			t.Fatalf("NewOTP() returned an error: %v", err)
		}

		if otp == nil {
			t.Fatal("NewOTP() returned a nil OTP, but no error.")
		}

		const otpLength = 6
		if len(otp.String()) != otpLength {
			t.Errorf("Expected OTP length to be %d, but got %d", otpLength, len(otp.String()))
		}

		// Check if the OTP consists only of digits
		match, err := regexp.MatchString("^[0-9]+$", otp.String())
		if err != nil {
			t.Fatalf("regex matching failed: %v", err)
		}

		if !match {
			t.Errorf("Expected OTP to contain only digits, but got %s", otp.String())
		}
	})
}

func TestFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid 6-digit code",
			input:   "123456",
			wantErr: false,
		},
		{
			name:    "valid code with leading zeros",
			input:   "000123",
			wantErr: false,
		},
		{
			name:    "invalid - contains non-digits",
			input:   "12345a",
			wantErr: true,
		},
		{
			name:    "invalid - too short (5 digits)",
			input:   "12345",
			wantErr: true,
		},
		{
			name:    "invalid - too long (7 digits)",
			input:   "1234567",
			wantErr: true,
		},
		{
			name:    "invalid - empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid - contains spaces",
			input:   "123 456",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := otp.FromString(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FromString(%q) expected error, but got nil", tt.input)
				}
				if result != nil {
					t.Errorf("FromString(%q) expected nil result on error, but got %v", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("FromString(%q) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("FromString(%q) expected non-nil result, but got nil", tt.input)
				}
				if result != nil && result.String() != tt.input {
					t.Errorf("FromString(%q) expected String() to return %q, but got %q", tt.input, tt.input, result.String())
				}
			}
		})
	}
}
