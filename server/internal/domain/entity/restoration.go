package entity

import (
	"errors"
	"time"

	"custom_auth_api/internal/domain/vo/email"
	"custom_auth_api/internal/domain/vo/ipaddress"
	"custom_auth_api/internal/domain/vo/otp"
)

// RestorationData contains all fields needed to restore a persisted OTPSession.
// This type is designed to be used exclusively by repository implementations
// when reconstructing sessions from persistent storage.
//
// Use NewRestorationData() to create instances with validation.
type RestorationData struct {
	Email         *email.Email
	Code          *otp.OTP
	Attempts      int
	CreatedAt     time.Time
	ExpiresAt     time.Time
	IPAddressHash *ipaddress.Hash
	UserAgent     string
}

// NewRestorationData creates restoration data with validation.
// This ensures that all required fields are present and valid before
// restoring an OTPSession from persistent storage.
//
// Validation rules:
//   - Email must not be nil
//   - OTP code must not be nil
//   - Attempts must not be negative
//   - CreatedAt must not be zero time
//   - ExpiresAt must not be zero time
//   - IPAddressHash must not be nil (use ipaddress.NewEmptyHash() if no IP)
//   - UserAgent can be empty string (optional field)
func NewRestorationData(
	userEmail *email.Email,
	otpCode *otp.OTP,
	attempts int,
	createdAt time.Time,
	expiresAt time.Time,
	ipHash *ipaddress.Hash,
	userAgent string,
) (*RestorationData, error) {
	if userEmail == nil {
		return nil, errors.New("email is required for restoration")
	}
	if otpCode == nil {
		return nil, errors.New("otp code is required for restoration")
	}
	if attempts < 0 {
		return nil, errors.New("attempts cannot be negative")
	}
	if createdAt.IsZero() {
		return nil, errors.New("createdAt is required for restoration")
	}
	if expiresAt.IsZero() {
		return nil, errors.New("expiresAt is required for restoration")
	}
	if ipHash == nil {
		return nil, errors.New("ipAddressHash is required for restoration (use NewEmptyHash if no IP)")
	}

	return &RestorationData{
		Email:         userEmail,
		Code:          otpCode,
		Attempts:      attempts,
		CreatedAt:     createdAt,
		ExpiresAt:     expiresAt,
		IPAddressHash: ipHash,
		UserAgent:     userAgent,
	}, nil
}

// RestoreOTPSession reconstructs an OTPSession from persisted data.
//
// REPOSITORY USE ONLY: This function is intended exclusively for repository
// implementations when loading sessions from persistent storage (Firestore, etc.).
// Application code should NEVER call this directly - use NewOTPSession instead.
//
// This function bypasses the normal constructor logic and directly sets all fields,
// including mutable state like the attempts counter, to their persisted values.
// This allows the session to resume from its last known state.
func RestoreOTPSession(data *RestorationData) *OTPSession {
	return &OTPSession{
		email:         data.Email,
		code:          data.Code,
		attempts:      data.Attempts,
		createdAt:     data.CreatedAt,
		expiresAt:     data.ExpiresAt,
		ipAddressHash: data.IPAddressHash,
		userAgent:     data.UserAgent,
	}
}
