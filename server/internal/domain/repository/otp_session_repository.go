package repository

import (
	"context"

	"custom_auth_api/internal/domain/entity"
	"custom_auth_api/internal/domain/vo/email"
)

// OTPSessionRepository defines the interface for OTPSession persistence.
// This is a clean repository interface with NO business logic.
// All business rules (expiration checks, attempt limits) are handled by the entity.
type OTPSessionRepository interface {
	// Save stores or updates an OTP session.
	// The session contains all necessary information including OTP code, expiration, attempts, etc.
	Save(ctx context.Context, session *entity.OTPSession) error

	// FindByEmail retrieves an OTP session by user email.
	// Returns entity.ErrSessionNotFound if no session exists for the email.
	// Does NOT perform business logic checks (expiration, attempts) - that's the entity's responsibility.
	FindByEmail(ctx context.Context, userEmail *email.Email) (*entity.OTPSession, error)

	// Delete removes an OTP session by user email.
	// Used after successful verification (one-time use) or for cleanup.
	Delete(ctx context.Context, userEmail *email.Email) error
}
