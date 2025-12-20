package otp

import (
	"crypto/rand"
	"fmt"
	"io"
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
