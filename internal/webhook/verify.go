package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"
)

// ValidateSignature reads the request body and validates the GitHub
// HMAC-SHA256 signature from the X-Hub-Signature-256 header.
// Returns the raw body bytes on success.
func ValidateSignature(r *http.Request, secret string) ([]byte, error) {
	if secret == "" {
		return nil, errors.New("webhook secret not configured")
	}

	sigHeader := r.Header.Get("X-Hub-Signature-256")
	if sigHeader == "" {
		return nil, errors.New("missing X-Hub-Signature-256 header")
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MB max
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(strings.TrimSpace(sigHeader))) {
		return nil, errors.New("signature mismatch")
	}

	return body, nil
}
