package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
)

// Secret returns a redacted string when displayed.
type Secret string

// String returns a redacted message to prevent the secret from being displayed.
func (secret Secret) String() string {
	return "[REDACTED]"
}

// Equals performs a constant time comparison to determine if the provided secret is equal.
func (secret Secret) Equals(value Secret) bool {
	if subtle.ConstantTimeCompare([]byte(secret), []byte(value)) == 1 {
		return true
	}
	return false
}

// NewSecret returns a new base64 encoded secret of the provided length.
func NewSecret(length int) Secret {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return Secret(base64.StdEncoding.EncodeToString(b))
}
