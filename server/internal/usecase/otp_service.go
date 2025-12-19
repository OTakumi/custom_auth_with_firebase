package usecase

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

const otpLength = 6

// OTPService handles OTP related business logic, such as generation, storage, and sending.
type OTPService struct{}

// NewOTPService creates a new OTPService.
func NewOTPService() *OTPService {
	return &OTPService{}
}

// GenerateOTP generates a 6-digit one-time password and logs it.
func (s *OTPService) GenerateOTP(email string) (string, error) {
	otp, err := generate6DigitCode()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// For now, just print the OTP to the console.
	log.Printf("OTP for %s: %s", email, otp)

	return otp, nil
}

// generate6DigitCode generates a random 6-digit string.
func generate6DigitCode() (string, error) {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, otpLength)

	n, err := io.ReadAtLeast(rand.Reader, b, otpLength)
	if n != otpLength {
		return "", fmt.Errorf("failed to read enough random bytes: %w", err)
	}

	for i := range b {
		b[i] = table[int(b[i])%len(table)]
	}

	return string(b), nil
}
