package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"

	"custom_auth_api/internal/domain/emailsender"
	"custom_auth_api/internal/domain/repository"
)

const otpLength = 6

// OTPService handles OTP related business logic, such as generation, storage, and sending.
type OTPService struct {
	otpRepo     repository.OTPRepository
	emailSender emailsender.EmailSender
}

// NewOTPService creates a new OTPService.
func NewOTPService(otpRepo repository.OTPRepository, emailSender emailsender.EmailSender) *OTPService {
	return &OTPService{otpRepo: otpRepo, emailSender: emailSender}
}

// GenerateAndSendOTP generates a 6-digit one-time password, saves it to the repository, and sends it via email.
func (s *OTPService) GenerateAndSendOTP(ctx context.Context, email string) (string, error) {
	otp, err := generate6DigitCode()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Save the OTP to the repository
	err = s.otpRepo.Save(ctx, email, otp)
	if err != nil {
		return "", fmt.Errorf("failed to save OTP: %w", err)
	}

	// Send the OTP via email
	err = s.emailSender.SendOTP(ctx, email, otp)
	if err != nil {
		return "", fmt.Errorf("failed to send OTP email: %w", err)
	}

	// For now, also log the OTP to the console for visibility (optional after email sending is fully implemented).
	log.Printf("OTP for %s: %s (sent via email sender)", email, otp)

	return otp, nil
}

// generate6DigitCode generates a random 6-digit string.
func generate6DigitCode() (string, error) {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

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
