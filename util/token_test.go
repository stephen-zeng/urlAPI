package util

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

func TestNewSessionTokenFormat(t *testing.T) {
	token, err := NewSessionToken()
	if err != nil {
		t.Fatalf("NewSessionToken: %v", err)
	}
	if len(token) != 2*SessionTokenBytes {
		t.Fatalf("token length = %d, want %d", len(token), 2*SessionTokenBytes)
	}
	raw, err := hex.DecodeString(token)
	if err != nil {
		t.Fatalf("token is not hex: %v", err)
	}
	if len(raw) != SessionTokenBytes {
		t.Fatalf("decoded length = %d, want %d", len(raw), SessionTokenBytes)
	}
}

func TestNewSessionTokenUsesFullEntropySource(t *testing.T) {
	orig := tokenRandSource
	t.Cleanup(func() { tokenRandSource = orig })

	// Every byte of the source must end up in the token unchanged.
	pattern := make([]byte, SessionTokenBytes)
	for i := range pattern {
		pattern[i] = byte(i * 7)
	}
	tokenRandSource = strings.NewReader(string(pattern))
	token, err := NewSessionToken()
	if err != nil {
		t.Fatalf("NewSessionToken: %v", err)
	}
	if want := hex.EncodeToString(pattern); token != want {
		t.Fatalf("token = %s, want %s", token, want)
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("entropy unavailable") }

func TestNewSessionTokenPropagatesRandFailure(t *testing.T) {
	orig := tokenRandSource
	t.Cleanup(func() { tokenRandSource = orig })

	tokenRandSource = failingReader{}
	token, err := NewSessionToken()
	if err == nil {
		t.Fatal("expected an error when the entropy source fails")
	}
	if token != "" {
		t.Fatalf("token = %q, want empty on failure", token)
	}

	// A short read must also fail rather than produce a weak token.
	tokenRandSource = strings.NewReader("short")
	if _, err := NewSessionToken(); err == nil {
		t.Fatal("expected an error on a short read")
	}
}
