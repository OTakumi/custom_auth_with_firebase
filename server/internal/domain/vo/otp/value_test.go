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
