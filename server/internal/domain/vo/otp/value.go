package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const otpLength = 6

// OTP represents a one-time password value object.
type OTP struct {
	value string
}

// NewOTP generates a new OTP.
func NewOTP() (*OTP, error) {
	otp, err := generate6DigitCode()
	if err != nil {
		return nil, err
	}

	return &OTP{value: otp}, nil
}

// String returns the string representation of the OTP.
func (o *OTP) String() string {
	return o.value
}

// generate6DigitCode generates a random 6-digit string without modulo bias.
// Uses crypto/rand.Int for unbiased random number generation.
func generate6DigitCode() (string, error) {
	// Generate a number between 0 and 999999 (inclusive)
	max := big.NewInt(1000000) // 10^6 = 1,000,000

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate random OTP: %w", err)
	}

	// Format as 6-digit string with leading zeros
	return fmt.Sprintf("%06d", n.Int64()), nil
}
