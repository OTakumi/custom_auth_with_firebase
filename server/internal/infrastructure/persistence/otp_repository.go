package persistence

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	otpCollection = "otps"
	otpExpiration = 5 * time.Minute // OTPs are valid for 5 minutes
)

var (
	ErrOTPNotFound = errors.New("otp not found")
	ErrOTPExpired  = errors.New("otp has expired")
)

// otpDocument represents the structure of an OTP document in Firestore.
type otpDocument struct {
	OTP       string    `firestore:"otp"`
	ExpiresAt time.Time `firestore:"expiresAt"`
}

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
	_, err := r.client.Collection(otpCollection).Doc(email).Set(ctx, otpDocument{
		OTP:       otp,
		ExpiresAt: time.Now().Add(otpExpiration),
	})
	if err != nil {
		return fmt.Errorf("failed to save otp to firestore: %w", err)
	}

	return nil
}

// Find retrieves the OTP for a given email from Firestore.
// It returns the OTP string and an error if not found, expired, or other issues.
func (r *OTPRepository) Find(ctx context.Context, email string) (string, error) {
	docRef := r.client.Collection(otpCollection).Doc(email)

	docSnap, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", fmt.Errorf("%w for email: %s", ErrOTPNotFound, email)
		}

		return "", fmt.Errorf("failed to get otp from firestore: %w", err)
	}

	var otpDoc otpDocument

	err = docSnap.DataTo(&otpDoc)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal otp document: %w", err)
	}

	// Delete the OTP immediately after retrieval to ensure one-time use.
	// This should ideally be part of a transaction if other operations depend on its existence.
	_, err = docRef.Delete(ctx)
	if err != nil {
		// Log the error but don't return it as the OTP was successfully retrieved.
		// The main goal is to return the OTP for verification.
		log.Printf("failed to delete otp for email %s: %v\n", email, err)
	}

	if time.Now().After(otpDoc.ExpiresAt) {
		return "", fmt.Errorf("%w for email: %s", ErrOTPExpired, email)
	}

	return otpDoc.OTP, nil
}
