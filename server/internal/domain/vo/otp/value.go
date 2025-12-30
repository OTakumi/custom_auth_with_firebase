package otp

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
)

// OTP represents a one-time password value object.
type OTP struct {
	value string
}

var (
	// ErrInvalidOTPFormat is returned when the OTP format is invalid.
	ErrInvalidOTPFormat = errors.New("otp must be exactly 6 digits")
	otpPattern          = regexp.MustCompile(`^[0-9]{6}$`)
)

// NewOTP generates a new OTP.
func NewOTP() (*OTP, error) {
	otp, err := generate6DigitCode()
	if err != nil {
		return nil, err
	}

	return &OTP{value: otp}, nil
}

// FromString creates an OTP from a string value.
// Returns an error if the string is not exactly 6 digits.
func FromString(code string) (*OTP, error) {
	if !otpPattern.MatchString(code) {
		return nil, ErrInvalidOTPFormat
	}

	return &OTP{value: code}, nil
}

// String returns the string representation of the OTP.
func (o *OTP) String() string {
	return o.value
}

// generate6DigitCode generates a random 6-digit string without modulo bias.
// Uses crypto/rand.Int for unbiased random number generation.
func generate6DigitCode() (string, error) {
	const otpMaxValue = 1000000 // 10^6 = 1,000,000
	// Generate a number between 0 and 999999 (inclusive)
	maxValue := big.NewInt(otpMaxValue)

	n, err := rand.Int(rand.Reader, maxValue)
	if err != nil {
		return "", fmt.Errorf("failed to generate random OTP: %w", err)
	}

	// Format as 6-digit string with leading zeros
	return fmt.Sprintf("%06d", n.Int64()), nil
}
