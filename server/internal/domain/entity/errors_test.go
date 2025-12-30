package entity_test

import "custom_auth_api/internal/domain/entity"

import (
	"errors"
	"testing"
)

func TestDomainErrors_AreDefined(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "entity.ErrSessionNotFound has correct message",
			err:  entity.ErrSessionNotFound,
			want: "otp session not found",
		},
		{
			name: "entity.ErrSessionExpired has correct message",
			err:  entity.ErrSessionExpired,
			want: "otp session has expired",
		},
		{
			name: "entity.ErrTooManyAttempts has correct message",
			err:  entity.ErrTooManyAttempts,
			want: "too many failed verification attempts",
		},
		{
			name: "entity.ErrInvalidOTP has correct message",
			err:  entity.ErrInvalidOTP,
			want: "invalid otp code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("expected error to be defined, got nil")
			}
			if tt.err.Error() != tt.want {
				t.Errorf("error message = %q, want %q", tt.err.Error(), tt.want)
			}
		})
	}
}

func TestDomainErrors_AreDistinct(t *testing.T) {
	// Ensure all errors are distinct from each other
	if errors.Is(entity.ErrSessionNotFound, entity.ErrSessionExpired) {
		t.Error("ErrSessionNotFound should not equal ErrSessionExpired")
	}
	if errors.Is(entity.ErrSessionNotFound, entity.ErrTooManyAttempts) {
		t.Error("ErrSessionNotFound should not equal ErrTooManyAttempts")
	}
	if errors.Is(entity.ErrSessionNotFound, entity.ErrInvalidOTP) {
		t.Error("ErrSessionNotFound should not equal ErrInvalidOTP")
	}
	if errors.Is(entity.ErrSessionExpired, entity.ErrTooManyAttempts) {
		t.Error("ErrSessionExpired should not equal ErrTooManyAttempts")
	}
	if errors.Is(entity.ErrSessionExpired, entity.ErrInvalidOTP) {
		t.Error("ErrSessionExpired should not equal ErrInvalidOTP")
	}
	if errors.Is(entity.ErrTooManyAttempts, entity.ErrInvalidOTP) {
		t.Error("ErrTooManyAttempts should not equal ErrInvalidOTP")
	}
}
