package usecase

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

// AuthService handles authentication related business logic, such as custom token generation.
type AuthService struct {
	authClient *auth.Client
}

// NewAuthService creates a new AuthService.
func NewAuthService(authClient *auth.Client) *AuthService {
	return &AuthService{authClient: authClient}
}

// GenerateCustomToken generates a custom Firebase authentication token for the given user ID (UID).
func (s *AuthService) GenerateCustomToken(ctx context.Context, uid string) (string, error) {
	customToken, err := s.authClient.CustomToken(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("failed to generate custom token: %w", err)
	}

	return customToken, nil
}
