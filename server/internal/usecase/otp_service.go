package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"custom_auth_api/internal/domain/emailsender"
	"custom_auth_api/internal/domain/repository"
	"custom_auth_api/internal/domain/vo/otp"
)

var (
	ErrInvalidOTP = errors.New("invalid OTP")
)

// OTPService handles OTP (One-Time Password) related business logic.
//
// Responsibilities:
// - Generate cryptographically secure 6-digit OTP codes
// - Store OTP with 5-minute expiration in the repository
// - Send OTP to users via email sender
// - Verify OTP against stored value
// - Track failed verification attempts (max 3 attempts)
// - Delete OTP after successful verification (one-time use)
//
// Security Features:
// - Attempt limiting: Maximum 3 failed verification attempts
// - Time-based expiration: OTP valid for 5 minutes
// - One-time use: OTP deleted after successful verification
// - Secure random generation: Uses crypto/rand without modulo bias
//
// Note:
// - User existence validation is handled by AuthService
// - Email format validation is handled at the handler layer.
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
	newOtp, err := otp.NewOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	otpStr := newOtp.String()

	// Save the OTP to the repository
	err = s.otpRepo.Save(ctx, email, otpStr)
	if err != nil {
		return "", fmt.Errorf("failed to save OTP: %w", err)
	}

	// Send the OTP via email
	err = s.emailSender.SendOTP(ctx, email, otpStr)
	if err != nil {
		return "", fmt.Errorf("failed to send OTP email: %w", err)
	}

	return otpStr, nil
}

// VerifyOTP validates the provided OTP against the stored one.
func (s *OTPService) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	storedOTP, err := s.otpRepo.Find(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve OTP for verification: %w", err)
	}

	if storedOTP != otp {
		// Increment failed attempts counter
		incrementErr := s.otpRepo.IncrementAttempts(ctx, email)
		if incrementErr != nil {
			log.Printf("Warning: failed to increment OTP attempts for %s: %v", email, incrementErr)
		}

		return false, fmt.Errorf("%w for email: %s", ErrInvalidOTP, email)
	}

	// Delete OTP after successful verification (one-time use)
	deleteErr := s.otpRepo.Delete(ctx, email)
	if deleteErr != nil {
		log.Printf("Warning: failed to delete OTP for %s: %v", email, deleteErr)
	}

	return true, nil
}
