package repository

import "context"

// OTPRepository defines the interface for OTP persistence.
type OTPRepository interface {
	Save(ctx context.Context, email, otp string) error
	Find(ctx context.Context, email string) (string, error)
}
