package entity

import (
	"crypto/subtle"
	"time"

	"custom_auth_api/internal/domain/vo/email"
	"custom_auth_api/internal/domain/vo/ipaddress"
	"custom_auth_api/internal/domain/vo/otp"
)

const (
	// DefaultOTPExpiration is the default duration for which an OTP session is valid.
	DefaultOTPExpiration = 5 * time.Minute

	// MaxVerificationAttempts is the maximum number of failed verification attempts allowed.
	MaxVerificationAttempts = 3
)

// OTPSession represents an OTP verification session for a user.
// This is an Entity (not a Value Object) because:
//   - It has identity (email)
//   - It has mutable state (attempts counter)
//   - It has lifecycle (created → verified/expired → deleted)
//
// Immutability: All fields except 'attempts' are immutable after creation.
// The attempts counter can only be modified through RecordFailedAttempt() method.
type OTPSession struct {
	email         *email.Email
	code          *otp.OTP
	createdAt     time.Time
	expiresAt     time.Time
	ipAddressHash *ipaddress.Hash // SHA-256 hash of IP address for privacy compliance
	userAgent     string
	attempts      int // Changes during verification attempts
}

// NewOTPSession creates a new OTP session with default expiration (5 minutes).
// All fields are immutable except the attempts counter.
func NewOTPSession(userEmail *email.Email, otpCode *otp.OTP) *OTPSession {
	now := time.Now()

	return &OTPSession{
		email:         userEmail,
		code:          otpCode,
		createdAt:     now,
		expiresAt:     now.Add(DefaultOTPExpiration),
		ipAddressHash: ipaddress.NewEmptyHash(),
		userAgent:     "",
		attempts:      0,
	}
}

// NewOTPSessionWithContext creates a new OTP session with IP and User-Agent for audit trail.
// IP addresses are hashed using SHA-256 for privacy compliance (GDPR).
func NewOTPSessionWithContext(
	userEmail *email.Email,
	otpCode *otp.OTP,
	ipAddress string,
	userAgent string,
) *OTPSession {
	session := NewOTPSession(userEmail, otpCode)
	session.ipAddressHash = ipaddress.NewHash(ipAddress)
	session.userAgent = userAgent

	return session
}

// Verify checks if the provided OTP code matches the stored code.
// Returns nil on successful verification.
// Returns ErrSessionExpired if the session has expired.
// Returns ErrTooManyAttempts if max attempts (3) have been exceeded.
// Returns ErrInvalidOTP if the code doesn't match.
//
// Uses constant-time comparison to prevent timing attacks.
// Automatically increments the attempts counter on mismatch.
func (s *OTPSession) Verify(inputCode string) error {
	// Check if session is eligible for verification
	err := s.CanVerify()
	if err != nil {
		return err
	}

	// Timing-safe comparison to prevent timing attacks
	expected := []byte(s.code.String())
	actual := []byte(inputCode)

	// Check length first (constant-time compare requires same length)
	if len(expected) != len(actual) {
		s.attempts++

		return ErrInvalidOTP
	}

	// Constant-time comparison
	if subtle.ConstantTimeCompare(expected, actual) != 1 {
		s.attempts++

		return ErrInvalidOTP
	}

	// Successful verification - do not increment attempts
	return nil
}

// CanVerify checks if the session is eligible for verification.
// Returns nil if the session can be verified.
// Returns ErrSessionExpired if expired.
// Returns ErrTooManyAttempts if max attempts exceeded.
func (s *OTPSession) CanVerify() error {
	if s.IsExpired() {
		return ErrSessionExpired
	}

	if s.attempts >= MaxVerificationAttempts {
		return ErrTooManyAttempts
	}

	return nil
}

// IsExpired checks if the OTP session has expired.
func (s *OTPSession) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

// RecordFailedAttempt increments the failed verification attempts counter.
// This is the only method that modifies the mutable state of the entity.
func (s *OTPSession) RecordFailedAttempt() {
	s.attempts++
}

// Getters for immutable fields

// Email returns the user's email address.
func (s *OTPSession) Email() *email.Email {
	return s.email
}

// OTP returns the OTP code (for repository serialization).
func (s *OTPSession) OTP() *otp.OTP {
	return s.code
}

// Attempts returns the current number of failed verification attempts.
func (s *OTPSession) Attempts() int {
	return s.attempts
}

// CreatedAt returns the session creation timestamp.
func (s *OTPSession) CreatedAt() time.Time {
	return s.createdAt
}

// ExpiresAt returns the session expiration timestamp.
func (s *OTPSession) ExpiresAt() time.Time {
	return s.expiresAt
}

// IPAddressHash returns the SHA-256 hash of the IP address.
func (s *OTPSession) IPAddressHash() *ipaddress.Hash {
	return s.ipAddressHash
}

// UserAgent returns the user agent string.
// Returns empty string if no user agent was provided.
func (s *OTPSession) UserAgent() string {
	return s.userAgent
}
