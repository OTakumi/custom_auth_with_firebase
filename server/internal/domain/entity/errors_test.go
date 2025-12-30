package entity

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
			name: "ErrSessionNotFound has correct message",
			err:  ErrSessionNotFound,
			want: "otp session not found",
		},
		{
			name: "ErrSessionExpired has correct message",
			err:  ErrSessionExpired,
			want: "otp session has expired",
		},
		{
			name: "ErrTooManyAttempts has correct message",
			err:  ErrTooManyAttempts,
			want: "too many failed verification attempts",
		},
		{
			name: "ErrInvalidOTP has correct message",
			err:  ErrInvalidOTP,
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
	if errors.Is(ErrSessionNotFound, ErrSessionExpired) {
		t.Error("ErrSessionNotFound should not equal ErrSessionExpired")
	}
	if errors.Is(ErrSessionNotFound, ErrTooManyAttempts) {
		t.Error("ErrSessionNotFound should not equal ErrTooManyAttempts")
	}
	if errors.Is(ErrSessionNotFound, ErrInvalidOTP) {
		t.Error("ErrSessionNotFound should not equal ErrInvalidOTP")
	}
	if errors.Is(ErrSessionExpired, ErrTooManyAttempts) {
		t.Error("ErrSessionExpired should not equal ErrTooManyAttempts")
	}
	if errors.Is(ErrSessionExpired, ErrInvalidOTP) {
		t.Error("ErrSessionExpired should not equal ErrInvalidOTP")
	}
	if errors.Is(ErrTooManyAttempts, ErrInvalidOTP) {
		t.Error("ErrTooManyAttempts should not equal ErrInvalidOTP")
	}
}
