package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	otpCollection  = "otps"
	otpExpiration  = 5 * time.Minute // OTPs are valid for 5 minutes
	maxOTPAttempts = 3               // Maximum OTP verification attempts
)

var (
	ErrOTPNotFound     = errors.New("otp not found")
	ErrOTPExpired      = errors.New("otp has expired")
	ErrTooManyAttempts = errors.New("too many failed attempts")
)

// otpDocument represents the structure of an OTP document in Firestore.
type otpDocument struct {
	OTP       string    `firestore:"otp"`
	ExpiresAt time.Time `firestore:"expiresAt"`
	Attempts  int       `firestore:"attempts"` // Track failed verification attempts
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
		Attempts:  0, // Initialize attempts to 0
	})
	if err != nil {
		return fmt.Errorf("failed to save otp to firestore: %w", err)
	}

	return nil
}

// Find retrieves the OTP for a given email from Firestore.
// It returns the OTP string and an error if not found, expired, or other issues.
// Does NOT delete the OTP - deletion must be done separately after successful verification.
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

	// Check if too many attempts have been made
	if otpDoc.Attempts >= maxOTPAttempts {
		return "", fmt.Errorf("%w for email: %s", ErrTooManyAttempts, email)
	}

	// Check expiration
	if time.Now().After(otpDoc.ExpiresAt) {
		return "", fmt.Errorf("%w for email: %s", ErrOTPExpired, email)
	}

	return otpDoc.OTP, nil
}

// Delete removes the OTP document for the given email.
func (r *OTPRepository) Delete(ctx context.Context, email string) error {
	_, err := r.client.Collection(otpCollection).Doc(email).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete otp: %w", err)
	}

	return nil
}

// IncrementAttempts increments the failed verification attempts counter for the given email.
func (r *OTPRepository) IncrementAttempts(ctx context.Context, email string) error {
	docRef := r.client.Collection(otpCollection).Doc(email)

	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "attempts", Value: firestore.Increment(1)},
	})
	if err != nil {
		// If document doesn't exist, it's already been deleted or never existed
		if status.Code(err) == codes.NotFound {
			return nil
		}
		return fmt.Errorf("failed to increment attempts: %w", err)
	}

	return nil
}
