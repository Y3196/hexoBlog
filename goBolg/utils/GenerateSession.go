// utils/session.go

package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSessionID creates a new session ID.
func GenerateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
