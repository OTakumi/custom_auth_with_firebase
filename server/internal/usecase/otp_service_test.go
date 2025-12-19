package usecase_test

import (
	"custom_auth_api/internal/usecase"
	"regexp"
	"testing"
)

func TestOTPService_GenerateOTP(t *testing.T) {
	t.Parallel()
	sut := usecase.NewOTPService()

	t.Run("should generate a 6-digit OTP", func(t *testing.T) {
		t.Parallel()
		otp, err := sut.GenerateOTP("test@example.com")

		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}

		if len(otp) != 6 {
			t.Errorf("expected OTP length to be 6, but got %d", len(otp))
		}

		isDigit := regexp.MustCompile(`^[0-9]+$`).MatchString
		if !isDigit(otp) {
			t.Errorf("expected OTP to be composed of digits, but got %s", otp)
		}
	})
}
