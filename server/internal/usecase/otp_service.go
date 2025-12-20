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

	// For now, also log the OTP to the console for visibility (optional after email sending is fully implemented).
	log.Printf("OTP for %s: %s (sent via email sender)", email, otpStr)

	return otpStr, nil
}

// VerifyOTP validates the provided OTP against the stored one.
func (s *OTPService) VerifyOTP(ctx context.Context, email, otp string) (bool, error) {
	storedOTP, err := s.otpRepo.Find(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve OTP for verification: %w", err)
	}

	if storedOTP != otp {
		return false, fmt.Errorf("%w for email: %s", ErrInvalidOTP, email)
	}

	return true, nil
}
