package httpserver

import (
	"crypto/rand"
	"encoding/hex"
)

const RequestIDHeader = "X-Request-ID"

func newRequestID() string {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(b[:])
}
