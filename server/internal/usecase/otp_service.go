package usecase

import (
	"context"
	"fmt"

	"custom_auth_api/internal/domain/emailsender"
	"custom_auth_api/internal/domain/entity"
	"custom_auth_api/internal/domain/repository"
	"custom_auth_api/internal/domain/vo/email"
	"custom_auth_api/internal/domain/vo/otp"
)

// OTPService handles OTP (One-Time Password) operations using entity-based design.
//
// Responsibilities:
// - Orchestrate OTP session creation and email delivery
// - Delegate business logic to OTPSession entity
// - Coordinate between repository and email sender
//
// Business Rules (delegated to OTPSession entity):
// - OTP expiration: 5 minutes (defined in entity.DefaultOTPExpiration)
// - Maximum verification attempts: 3 (defined in entity.MaxVerificationAttempts)
// - One-time use: Session deleted after successful verification
// - Timing-safe comparison for OTP verification
//
// Note:
// - User existence validation is handled by AuthService
// - Email format validation is handled by email value object
type OTPService struct {
	sessionRepo repository.OTPSessionRepository
	emailSender emailsender.EmailSender
}

// NewOTPService creates a new OTPService.
func NewOTPService(sessionRepo repository.OTPSessionRepository, emailSender emailsender.EmailSender) *OTPService {
	return &OTPService{
		sessionRepo: sessionRepo,
		emailSender: emailSender,
	}
}

// GenerateAndSendOTP generates a new OTP session and sends the OTP code via email.
// Returns the generated OTP code string (for testing purposes).
func (s *OTPService) GenerateAndSendOTP(ctx context.Context, emailAddr string) (string, error) {
	// Validate and create email value object
	userEmail, err := email.NewEmail(emailAddr)
	if err != nil {
		return "", fmt.Errorf("invalid email address: %w", err)
	}

	// Generate OTP code
	otpCode, err := otp.NewOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Create new OTP session entity
	session := entity.NewOTPSession(userEmail, otpCode)

	// Persist the session
	err = s.sessionRepo.Save(ctx, session)
	if err != nil {
		return "", fmt.Errorf("failed to save OTP session: %w", err)
	}

	// Send OTP via email
	err = s.emailSender.SendOTP(ctx, emailAddr, otpCode.String())
	if err != nil {
		return "", fmt.Errorf("failed to send OTP email: %w", err)
	}

	return otpCode.String(), nil
}

// VerifyOTP validates the provided OTP code against the stored session.
// Returns true if verification succeeds, false otherwise.
// Automatically handles:
// - Expiration checking (via entity)
// - Attempt counting (via entity)
// - Session deletion on success.
func (s *OTPService) VerifyOTP(ctx context.Context, emailAddr, inputCode string) (bool, error) {
	// Validate and create email value object
	userEmail, err := email.NewEmail(emailAddr)
	if err != nil {
		return false, fmt.Errorf("invalid email address: %w", err)
	}

	// Retrieve session from repository
	session, err := s.sessionRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve OTP session: %w", err)
	}

	// Verify OTP using entity business logic
	err = session.Verify(inputCode)
	if err != nil {
		// Save updated attempts count (entity incremented it on failure)
		saveErr := s.sessionRepo.Save(ctx, session)
		if saveErr != nil {
			// Log warning but return the original verification error
			return false, fmt.Errorf("verification failed (%w) and save failed: %w", err, saveErr)
		}

		return false, err
	}

	// Successful verification - delete session (one-time use)
	err = s.sessionRepo.Delete(ctx, userEmail)
	if err != nil {
		// Verification succeeded but cleanup failed
		// The session will eventually expire naturally
		return true, fmt.Errorf("OTP verified but cleanup failed: %w", err)
	}

	return true, nil
}
