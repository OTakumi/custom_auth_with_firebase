package ipaddress_test

import "custom_auth_api/internal/domain/vo/ipaddress"

import (
	"testing"
)

const (
	testIP = "192.168.1.1"
)

func TestNewHash(t *testing.T) {
	t.Parallel()

	t.Run("creates hash from IP address", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ip := testIP

		// Act
		hash := ipaddress.NewHash(ip)

		// Assert
		if hash == nil {
			t.Fatal("expected non-nil hash")
		}
		if hash.String() == "" {
			t.Error("expected non-empty hash value")
		}
		// SHA-256 produces 64-character hex string
		if len(hash.String()) != 64 {
			t.Errorf("expected hash length 64, got %d", len(hash.String()))
		}
	})

	t.Run("same IP produces same hash (deterministic)", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ip := testIP

		// Act
		hash1 := ipaddress.NewHash(ip)
		hash2 := ipaddress.NewHash(ip)

		// Assert
		if hash1.String() != hash2.String() {
			t.Error("same IP should produce same hash")
		}
	})

	t.Run("different IPs produce different hashes", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ip1 := testIP
		ip2 := "192.168.1.2"

		// Act
		hash1 := ipaddress.NewHash(ip1)
		hash2 := ipaddress.NewHash(ip2)

		// Assert
		if hash1.String() == hash2.String() {
			t.Error("different IPs should produce different hashes")
		}
	})

	t.Run("handles empty string", func(t *testing.T) {
		t.Parallel()

		// Act
		hash := ipaddress.NewHash("")

		// Assert
		if hash == nil {
			t.Fatal("expected non-nil hash")
		}
		// Empty string still produces a hash
		if hash.String() == "" {
			t.Error("expected non-empty hash even for empty input")
		}
		if len(hash.String()) != 64 {
			t.Errorf("expected hash length 64, got %d", len(hash.String()))
		}
	})

	t.Run("produces valid hex characters only", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ip := "10.0.0.1"

		// Act
		hash := ipaddress.NewHash(ip)

		// Assert
		for _, c := range hash.String() {
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
				t.Errorf("hash contains invalid hex character: %c", c)
			}
		}
	})
}

func TestNewEmptyHash(t *testing.T) {
	t.Parallel()

	t.Run("creates empty hash", func(t *testing.T) {
		t.Parallel()

		// Act
		hash := ipaddress.NewEmptyHash()

		// Assert
		if hash == nil {
			t.Fatal("expected non-nil hash")
		}
		if hash.String() != "" {
			t.Errorf("expected empty string, got %q", hash.String())
		}
		if !hash.IsEmpty() {
			t.Error("IsEmpty() should return true for empty hash")
		}
	})
}

func TestHash_IsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		hash     *ipaddress.Hash
		expected bool
	}{
		{
			name:     "empty hash returns true",
			hash:     ipaddress.NewEmptyHash(),
			expected: true,
		},
		{
			name:     "non-empty hash returns false",
			hash:     ipaddress.NewHash(testIP),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			result := tt.hash.IsEmpty()

			// Assert
			if result != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHash_String(t *testing.T) {
	t.Parallel()

	t.Run("returns the hash value", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ip := "203.0.113.0"
		hash := ipaddress.NewHash(ip)

		// Act
		value := hash.String()

		// Assert
		if value == "" {
			t.Error("expected non-empty string")
		}
		if len(value) != 64 {
			t.Errorf("expected length 64, got %d", len(value))
		}
	})

	t.Run("empty hash returns empty string", func(t *testing.T) {
		t.Parallel()

		// Arrange
		hash := ipaddress.NewEmptyHash()

		// Act
		value := hash.String()

		// Assert
		if value != "" {
			t.Errorf("expected empty string, got %q", value)
		}
	})
}
