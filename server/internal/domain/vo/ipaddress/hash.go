package ipaddress

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash represents a SHA-256 hashed IP address for privacy protection.
// This value object ensures IP addresses are stored securely while still
// allowing fraud detection (same IP → same hash).
type Hash struct {
	value string
}

// NewHash creates a hash from a raw IP address string.
// Uses SHA-256 to one-way hash the IP address for GDPR/privacy compliance.
func NewHash(ipAddress string) *Hash {
	hash := sha256.Sum256([]byte(ipAddress))
	return &Hash{value: hex.EncodeToString(hash[:])}
}

// NewEmptyHash creates an empty hash for sessions without IP tracking.
func NewEmptyHash() *Hash {
	return &Hash{value: ""}
}

// FromString reconstructs a Hash from a previously computed hash string.
// This is used by repository implementations when loading from persistent storage.
// The input should be a 64-character hex string (SHA-256 hash).
func FromString(hashValue string) *Hash {
	return &Hash{value: hashValue}
}

// String returns the hash value (64-character hex string for SHA-256).
func (h *Hash) String() string {
	return h.value
}

// IsEmpty checks if the hash is empty.
func (h *Hash) IsEmpty() bool {
	return h.value == ""
}
