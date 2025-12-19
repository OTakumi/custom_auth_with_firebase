package persistence

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

const (
	otpCollection = "otps"
	otpExpiration = 5 * time.Minute // OTPs are valid for 5 minutes
)

// OTPRepository handles storing and retrieving OTPs from Firestore.
type OTPRepository struct {
	client *firestore.Client
}

// NewOTPRepository creates a new OTPRepository.
func NewOTPRepository(client *firestore.Client) *OTPRepository {
	return &OTPRepository{client: client}
}

// Save stores the OTP for a given email with an expiration time.
func (r *OTPRepository) Save(ctx context.Context, email, otp string) error {
	_, err := r.client.Collection(otpCollection).Doc(email).Set(ctx, map[string]any{
		"otp":       otp,
		"expiresAt": time.Now().Add(otpExpiration),
	})
	if err != nil {
		return fmt.Errorf("failed to save otp to firestore: %w", err)
	}

	return nil
}
