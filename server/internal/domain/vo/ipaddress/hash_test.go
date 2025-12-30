package ipaddress

import (
	"testing"
)

func TestNewHash(t *testing.T) {
	t.Parallel()

	t.Run("creates hash from IP address", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ip := "192.168.1.1"

		// Act
		hash := NewHash(ip)

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
		ip := "192.168.1.1"

		// Act
		hash1 := NewHash(ip)
		hash2 := NewHash(ip)

		// Assert
		if hash1.String() != hash2.String() {
			t.Error("same IP should produce same hash")
		}
	})

	t.Run("different IPs produce different hashes", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ip1 := "192.168.1.1"
		ip2 := "192.168.1.2"

		// Act
		hash1 := NewHash(ip1)
		hash2 := NewHash(ip2)

		// Assert
		if hash1.String() == hash2.String() {
			t.Error("different IPs should produce different hashes")
		}
	})

	t.Run("handles empty string", func(t *testing.T) {
		t.Parallel()

		// Act
		hash := NewHash("")

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
		hash := NewHash(ip)

		// Assert
		for _, c := range hash.String() {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
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
		hash := NewEmptyHash()

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
		hash     *Hash
		expected bool
	}{
		{
			name:     "empty hash returns true",
			hash:     NewEmptyHash(),
			expected: true,
		},
		{
			name:     "non-empty hash returns false",
			hash:     NewHash("192.168.1.1"),
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
		hash := NewHash(ip)

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
		hash := NewEmptyHash()

		// Act
		value := hash.String()

		// Assert
		if value != "" {
			t.Errorf("expected empty string, got %q", value)
		}
	})
}
