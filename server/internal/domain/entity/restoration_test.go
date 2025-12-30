package entity

import (
	"testing"
	"time"

	"custom_auth_api/internal/domain/vo/email"
	"custom_auth_api/internal/domain/vo/ipaddress"
	"custom_auth_api/internal/domain/vo/otp"
)

func TestNewRestorationData_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	userEmail, _ := email.NewEmail("test@example.com")
	otpCode, _ := otp.FromString("123456")
	attempts := 2
	createdAt := time.Now().Add(-4 * time.Minute)
	expiresAt := time.Now().Add(1 * time.Minute)
	ipHash := ipaddress.FromString("abc123hash")
	userAgent := "Mozilla/5.0"

	// Act
	data, err := NewRestorationData(
		userEmail,
		otpCode,
		attempts,
		createdAt,
		expiresAt,
		ipHash,
		userAgent,
	)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data.Email != userEmail {
		t.Error("email not set correctly")
	}
	if data.Code != otpCode {
		t.Error("code not set correctly")
	}
	if data.Attempts != attempts {
		t.Errorf("expected attempts %d, got %d", attempts, data.Attempts)
	}
	if !data.CreatedAt.Equal(createdAt) {
		t.Error("createdAt not set correctly")
	}
	if !data.ExpiresAt.Equal(expiresAt) {
		t.Error("expiresAt not set correctly")
	}
	if data.IPAddressHash.String() != ipHash.String() {
		t.Errorf("expected IP hash %q, got %q", ipHash, data.IPAddressHash)
	}
	if data.UserAgent != userAgent {
		t.Errorf("expected user agent %q, got %q", userAgent, data.UserAgent)
	}
}

func TestNewRestorationData_Validation(t *testing.T) {
	t.Parallel()

	validEmail, _ := email.NewEmail("test@example.com")
	validOTP, _ := otp.FromString("123456")
	validCreatedAt := time.Now().Add(-4 * time.Minute)
	validExpiresAt := time.Now().Add(1 * time.Minute)

	tests := []struct {
		name      string
		email     *email.Email
		code      *otp.OTP
		attempts  int
		createdAt time.Time
		expiresAt time.Time
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "nil email returns error",
			email:     nil,
			code:      validOTP,
			attempts:  0,
			createdAt: validCreatedAt,
			expiresAt: validExpiresAt,
			wantErr:   true,
			errMsg:    "email is required for restoration",
		},
		{
			name:      "nil OTP code returns error",
			email:     validEmail,
			code:      nil,
			attempts:  0,
			createdAt: validCreatedAt,
			expiresAt: validExpiresAt,
			wantErr:   true,
			errMsg:    "otp code is required for restoration",
		},
		{
			name:      "negative attempts returns error",
			email:     validEmail,
			code:      validOTP,
			attempts:  -1,
			createdAt: validCreatedAt,
			expiresAt: validExpiresAt,
			wantErr:   true,
			errMsg:    "attempts cannot be negative",
		},
		{
			name:      "zero createdAt returns error",
			email:     validEmail,
			code:      validOTP,
			attempts:  0,
			createdAt: time.Time{},
			expiresAt: validExpiresAt,
			wantErr:   true,
			errMsg:    "createdAt is required for restoration",
		},
		{
			name:      "zero expiresAt returns error",
			email:     validEmail,
			code:      validOTP,
			attempts:  0,
			createdAt: validCreatedAt,
			expiresAt: time.Time{},
			wantErr:   true,
			errMsg:    "expiresAt is required for restoration",
		},
		{
			name:      "zero attempts is valid",
			email:     validEmail,
			code:      validOTP,
			attempts:  0,
			createdAt: validCreatedAt,
			expiresAt: validExpiresAt,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			data, err := NewRestorationData(
				tt.email,
				tt.code,
				tt.attempts,
				tt.createdAt,
				tt.expiresAt,
				ipaddress.NewEmptyHash(),
				"",
			)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
				if data != nil {
					t.Error("expected nil data on error")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if data == nil {
					t.Error("expected non-nil data")
				}
			}
		})
	}
}

func TestRestoreOTPSession(t *testing.T) {
	t.Parallel()

	t.Run("restores session with all fields", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userEmail, _ := email.NewEmail("test@example.com")
		otpCode, _ := otp.FromString("123456")
		attempts := 2
		createdAt := time.Now().Add(-4 * time.Minute)
		expiresAt := time.Now().Add(1 * time.Minute)
		ipHash := ipaddress.FromString("abc123hash")
		userAgent := "Mozilla/5.0"

		data, _ := NewRestorationData(
			userEmail,
			otpCode,
			attempts,
			createdAt,
			expiresAt,
			ipHash,
			userAgent,
		)

		// Act
		session := RestoreOTPSession(data)

		// Assert
		if session == nil {
			t.Fatal("expected non-nil session")
		}
		if session.Email() != userEmail {
			t.Error("email not restored correctly")
		}
		if session.OTP() != otpCode {
			t.Error("OTP not restored correctly")
		}
		if session.Attempts() != attempts {
			t.Errorf("expected attempts %d, got %d", attempts, session.Attempts())
		}
		if !session.CreatedAt().Equal(createdAt) {
			t.Error("createdAt not restored correctly")
		}
		if !session.ExpiresAt().Equal(expiresAt) {
			t.Error("expiresAt not restored correctly")
		}
		if session.IPAddressHash() != ipHash {
			t.Errorf("expected IP hash %q, got %q", ipHash, session.IPAddressHash())
		}
		if session.UserAgent() != userAgent {
			t.Errorf("expected user agent %q, got %q", userAgent, session.UserAgent())
		}
	})

	t.Run("restored session maintains behavior", func(t *testing.T) {
		t.Parallel()

		// Arrange - create a session with 2 failed attempts
		userEmail, _ := email.NewEmail("test@example.com")
		otpCode, _ := otp.FromString("123456")
		data, _ := NewRestorationData(
			userEmail,
			otpCode,
			2, // 2 failed attempts already
			time.Now().Add(-4*time.Minute),
			time.Now().Add(1*time.Minute),
			ipaddress.FromString("hash"),
			"agent",
		)

		session := RestoreOTPSession(data)

		// Act - verify with wrong code (3rd attempt)
		err := session.Verify("wrong1")

		// Assert - should increment to 3
		if err != ErrInvalidOTP {
			t.Errorf("expected ErrInvalidOTP, got %v", err)
		}
		if session.Attempts() != 3 {
			t.Errorf("expected 3 attempts, got %d", session.Attempts())
		}

		// Act - 4th attempt should fail with too many attempts
		err = session.Verify("123456") // Even correct code should fail
		if err != ErrTooManyAttempts {
			t.Errorf("expected ErrTooManyAttempts, got %v", err)
		}
	})

	t.Run("restored expired session is expired", func(t *testing.T) {
		t.Parallel()

		// Arrange - create an expired session
		userEmail, _ := email.NewEmail("test@example.com")
		otpCode, _ := otp.FromString("123456")
		data, _ := NewRestorationData(
			userEmail,
			otpCode,
			0,
			time.Now().Add(-10*time.Minute), // created 10 min ago
			time.Now().Add(-5*time.Minute),  // expired 5 min ago
			ipaddress.FromString("hash"),
			"agent",
		)

		session := RestoreOTPSession(data)

		// Act & Assert
		if !session.IsExpired() {
			t.Error("expected session to be expired")
		}

		err := session.Verify("123456")
		if err != ErrSessionExpired {
			t.Errorf("expected ErrSessionExpired, got %v", err)
		}
	})
}
