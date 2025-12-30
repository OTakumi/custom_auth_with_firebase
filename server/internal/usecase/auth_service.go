package usecase

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

// AuthService handles Firebase Authentication related business logic.
//
// Responsibilities:
// - Retrieve user information from Firebase Auth
// - Generate Firebase custom tokens for authenticated users
//
// Note:
// - OTP generation, sending, and verification are handled by OTPService
// - This service acts as a wrapper around the Firebase Auth SDK.
type AuthService struct {
	authClient *auth.Client
}

// NewAuthService creates a new AuthService.
func NewAuthService(authClient *auth.Client) *AuthService {
	return &AuthService{authClient: authClient}
}

// GetUserByEmail retrieves a user by email address.
// Returns the user record if found, or an error if the user does not exist.
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	user, err := s.authClient.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GenerateCustomToken generates a custom Firebase authentication token for the given user ID (UID).
func (s *AuthService) GenerateCustomToken(ctx context.Context, uid string) (string, error) {
	customToken, err := s.authClient.CustomToken(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("failed to generate custom token: %w", err)
	}

	return customToken, nil
}
