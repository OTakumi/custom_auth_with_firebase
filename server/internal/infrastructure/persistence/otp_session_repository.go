package persistence

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"custom_auth_api/internal/domain/entity"
	"custom_auth_api/internal/domain/vo/email"
	"custom_auth_api/internal/domain/vo/ipaddress"
	"custom_auth_api/internal/domain/vo/otp"
)

const (
	otpSessionCollection = "otps"
)

// otpSessionDocument represents the Firestore document schema for OTP sessions.
// This is the persistence model, separate from the domain entity.
type otpSessionDocument struct {
	Email         string    `firestore:"email"`
	OTP           string    `firestore:"otp"`
	Attempts      int       `firestore:"attempts"`
	CreatedAt     time.Time `firestore:"createdAt"`
	ExpiresAt     time.Time `firestore:"expiresAt"`
	IPAddressHash string    `firestore:"ipAddressHash,omitempty"`
	UserAgent     string    `firestore:"userAgent,omitempty"`
}

// OTPSessionRepository handles OTPSession persistence in Firestore.
// This implementation contains NO business logic - it's purely for data access.
type OTPSessionRepository struct {
	client *firestore.Client
}

// NewOTPSessionRepository creates a new OTPSessionRepository.
func NewOTPSessionRepository(client *firestore.Client) *OTPSessionRepository {
	return &OTPSessionRepository{client: client}
}

// Save stores or updates an OTP session in Firestore.
// Uses the email as the document ID for deterministic lookups.
func (r *OTPSessionRepository) Save(ctx context.Context, session *entity.OTPSession) error {
	doc := otpSessionDocument{
		Email:         session.Email().Value,
		OTP:           session.OTP().String(),
		Attempts:      session.Attempts(),
		CreatedAt:     session.CreatedAt(),
		ExpiresAt:     session.ExpiresAt(),
		IPAddressHash: session.IPAddressHash().String(),
		UserAgent:     session.UserAgent(),
	}

	_, err := r.client.Collection(otpSessionCollection).Doc(session.Email().Value).Set(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to save otp session: %w", err)
	}

	return nil
}

// FindByEmail retrieves an OTP session from Firestore by email.
// Returns entity.ErrSessionNotFound if the document doesn't exist.
// Does NOT check expiration or attempt limits - that's the entity's responsibility.
func (r *OTPSessionRepository) FindByEmail(ctx context.Context, userEmail *email.Email) (*entity.OTPSession, error) {
	docSnap, err := r.client.Collection(otpSessionCollection).Doc(userEmail.Value).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, entity.ErrSessionNotFound
		}

		return nil, fmt.Errorf("failed to get otp session: %w", err)
	}

	var doc otpSessionDocument
	err = docSnap.DataTo(&doc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal otp session: %w", err)
	}

	// Reconstruct domain entity from persistence model
	otpCode, err := otp.FromString(doc.OTP)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct otp code: %w", err)
	}

	// Reconstruct email value object
	reconstructedEmail, err := email.NewEmail(doc.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct email: %w", err)
	}

	// Reconstruct the session from persisted data using validated RestorationData
	return reconstructSessionFromDocument(doc, reconstructedEmail, otpCode)
}

// Delete removes an OTP session from Firestore.
func (r *OTPSessionRepository) Delete(ctx context.Context, userEmail *email.Email) error {
	_, err := r.client.Collection(otpSessionCollection).Doc(userEmail.Value).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete otp session: %w", err)
	}

	return nil
}

// reconstructSessionFromDocument creates a domain entity from a Firestore document.
// Uses RestorationData to ensure type-safe reconstruction with validation.
func reconstructSessionFromDocument(
	doc otpSessionDocument,
	userEmail *email.Email,
	otpCode *otp.OTP,
) (*entity.OTPSession, error) {
	// Reconstruct IP address hash value object from stored hash string
	var ipHash *ipaddress.Hash
	if doc.IPAddressHash == "" {
		ipHash = ipaddress.NewEmptyHash()
	} else {
		ipHash = ipaddress.FromString(doc.IPAddressHash)
	}

	// Create validated restoration data from persisted fields
	restorationData, err := entity.NewRestorationData(
		userEmail,
		otpCode,
		doc.Attempts,
		doc.CreatedAt,
		doc.ExpiresAt,
		ipHash,
		doc.UserAgent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create restoration data: %w", err)
	}

	// Restore the session entity with all persisted state
	return entity.RestoreOTPSession(restorationData), nil
}
