package entity

import "errors"

var (
	// ErrSessionNotFound is returned when an OTP session is not found.
	ErrSessionNotFound = errors.New("otp session not found")

	// ErrSessionExpired is returned when an OTP session has expired.
	ErrSessionExpired = errors.New("otp session has expired")

	// ErrTooManyAttempts is returned when too many verification attempts have been made.
	ErrTooManyAttempts = errors.New("too many failed verification attempts")

	// ErrInvalidOTP is returned when the provided OTP does not match.
	ErrInvalidOTP = errors.New("invalid otp code")
)
