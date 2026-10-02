package util

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/pkg/errors"
)

// SessionTokenBytes is the amount of randomness in a session token (256 bits).
const SessionTokenBytes = 32

// tokenRandSource is the entropy source for tokens; tests may replace it.
var tokenRandSource io.Reader = rand.Reader

// NewSessionToken returns a hex-encoded token carrying 256 bits read directly
// from crypto/rand. The 64-character output keeps the same shape as the
// tokens issued by earlier versions, so stored sessions remain valid.
func NewSessionToken() (string, error) {
	buf := make([]byte, SessionTokenBytes)
	if _, err := io.ReadFull(tokenRandSource, buf); err != nil {
		return "", errors.Wrap(err, "generate session token")
	}
	return hex.EncodeToString(buf), nil
}
