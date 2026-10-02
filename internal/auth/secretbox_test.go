package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testKey(fill byte) []byte { return bytes.Repeat([]byte{fill}, 32) }

func TestSecretBoxRoundTrip(t *testing.T) {
	box, err := NewSecretBox(testKey(1))
	if err != nil {
		t.Fatal(err)
	}
	const secret = "sk-live-abcdef0123456789"
	sealed, err := box.Seal(secret, "provider/openai/api_key")
	if err != nil {
		t.Fatal(err)
	}
	if !IsEncrypted(sealed) || strings.Contains(sealed, secret) || strings.Contains(sealed, base64.StdEncoding.EncodeToString([]byte(secret))) {
		t.Fatalf("sealed value leaks secret: %s", sealed)
	}
	again, _ := box.Seal(secret, "provider/openai/api_key")
	if again == sealed {
		t.Fatal("nonce reuse: identical ciphertexts")
	}
	got, err := box.Open(sealed, "provider/openai/api_key")
	if err != nil || got != secret {
		t.Fatalf("Open = %q, %v", got, err)
	}
}

func TestSecretBoxRejectsWrongKeyContextAndTampering(t *testing.T) {
	box, _ := NewSecretBox(testKey(1))
	other, _ := NewSecretBox(testKey(2))
	sealed, _ := box.Seal("secret", "ctx")

	if _, err := other.Open(sealed, "ctx"); err == nil {
		t.Fatal("wrong key accepted")
	} else if strings.Contains(err.Error(), "secret\"") {
		t.Fatal("error leaks plaintext")
	}
	if _, err := box.Open(sealed, "other-ctx"); err == nil {
		t.Fatal("ciphertext accepted under a different context")
	}
	raw, _ := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(sealed, encryptedPrefix))
	raw[len(raw)-1] ^= 1
	if _, err := box.Open(encryptedPrefix+base64.RawURLEncoding.EncodeToString(raw), "ctx"); err == nil {
		t.Fatal("tampered ciphertext accepted")
	}
	for _, bad := range []string{"", "plain", "enc:v1:", "enc:v1:!!!", "enc:v1:AAAA"} {
		if _, err := box.Open(bad, "ctx"); err == nil {
			t.Errorf("Open(%q) accepted", bad)
		}
	}
}

func TestParseKey(t *testing.T) {
	key := testKey(7)
	for _, encoded := range []string{
		hex.EncodeToString(key),
		base64.StdEncoding.EncodeToString(key),
		base64.RawURLEncoding.EncodeToString(key),
		"  " + base64.StdEncoding.EncodeToString(key) + "\n",
	} {
		got, err := ParseKey(encoded)
		if err != nil || !bytes.Equal(got, key) {
			t.Errorf("ParseKey(%q) = %x, %v", encoded, got, err)
		}
	}
	for _, bad := range []string{"", "short", base64.StdEncoding.EncodeToString(key[:16]), strings.Repeat("z", 64)} {
		if _, err := ParseKey(bad); err == nil {
			t.Errorf("ParseKey(%q) accepted", bad)
		}
	}
	if _, err := NewSecretBox(key[:31]); err == nil {
		t.Fatal("short key accepted")
	}
}

func TestLoadSecretBoxFromEnv(t *testing.T) {
	t.Setenv(EnvSecretKey, "")
	t.Setenv(EnvSecretKeyFile, "")
	if box, err := LoadSecretBoxFromEnv(); box != nil || err != nil {
		t.Fatalf("unset: %v, %v", box, err)
	}

	t.Setenv(EnvSecretKey, "not-a-key")
	if _, err := LoadSecretBoxFromEnv(); err == nil || strings.Contains(err.Error(), "not-a-key") {
		t.Fatalf("invalid key: %v", err)
	}

	t.Setenv(EnvSecretKey, base64.StdEncoding.EncodeToString(testKey(3)))
	if box, err := LoadSecretBoxFromEnv(); box == nil || err != nil {
		t.Fatalf("valid key: %v, %v", box, err)
	}

	t.Setenv(EnvSecretKey, "")
	path := filepath.Join(t.TempDir(), "key")
	if err := os.WriteFile(path, []byte(hex.EncodeToString(testKey(4))+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvSecretKeyFile, path)
	if box, err := LoadSecretBoxFromEnv(); box == nil || err != nil {
		t.Fatalf("key file: %v, %v", box, err)
	}
	t.Setenv(EnvSecretKeyFile, filepath.Join(t.TempDir(), "missing"))
	if _, err := LoadSecretBoxFromEnv(); err == nil {
		t.Fatal("missing key file accepted")
	}
}
