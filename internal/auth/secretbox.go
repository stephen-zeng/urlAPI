package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const (
	// EnvSecretKey holds the 32-byte master key (base64 or hex) used to
	// encrypt stored provider secrets.
	EnvSecretKey = "URLAPI_SECRET_KEY"
	// EnvSecretKeyFile names a file containing the master key, for Docker
	// or systemd secrets. EnvSecretKey takes precedence.
	EnvSecretKeyFile = "URLAPI_SECRET_KEY_FILE"

	encryptedPrefix = "enc:v1:"
	masterKeyLen    = 32
)

// ErrNoSecretKey is returned when an encrypted value is read without a key.
var ErrNoSecretKey = errors.New("secret is encrypted but " + EnvSecretKey + " is not configured")

// SecretBox encrypts secrets with AES-256-GCM. Each value is bound to a
// context string (e.g. "provider/openai/api_key") as additional data, so
// ciphertexts cannot be swapped between fields.
type SecretBox struct {
	aead cipher.AEAD
}

// NewSecretBox creates a SecretBox from a 32-byte key.
func NewSecretBox(key []byte) (*SecretBox, error) {
	if len(key) != masterKeyLen {
		return nil, errors.Errorf("secret key must be %d bytes, got %d", masterKeyLen, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &SecretBox{aead: aead}, nil
}

// ParseKey decodes a 32-byte key given as hex or (raw/padded, standard or
// URL-safe) base64.
func ParseKey(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if len(s) == 2*masterKeyLen {
		if key, err := hex.DecodeString(s); err == nil {
			return key, nil
		}
	}
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
		if key, err := enc.DecodeString(s); err == nil {
			if len(key) != masterKeyLen {
				return nil, errors.Errorf("secret key must decode to %d bytes, got %d", masterKeyLen, len(key))
			}
			return key, nil
		}
	}
	return nil, errors.New("secret key must be 32 bytes encoded as base64 or hex")
}

// LoadSecretBoxFromEnv builds a SecretBox from EnvSecretKey or
// EnvSecretKeyFile. It returns (nil, nil) when neither is set.
func LoadSecretBoxFromEnv() (*SecretBox, error) {
	raw := os.Getenv(EnvSecretKey)
	source := EnvSecretKey
	if raw == "" {
		path := os.Getenv(EnvSecretKeyFile)
		if path == "" {
			return nil, nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "read %s", EnvSecretKeyFile)
		}
		raw, source = string(data), EnvSecretKeyFile
	}
	key, err := ParseKey(raw)
	if err != nil {
		return nil, errors.Wrap(err, source)
	}
	return NewSecretBox(key)
}

// IsEncrypted reports whether stored was produced by Seal.
func IsEncrypted(stored string) bool {
	return strings.HasPrefix(stored, encryptedPrefix)
}

// Seal encrypts plaintext for the given context.
func (b *SecretBox) Seal(plaintext, context string) (string, error) {
	nonce := make([]byte, b.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", errors.Wrap(err, "generate nonce")
	}
	sealed := b.aead.Seal(nonce, nonce, []byte(plaintext), []byte(context))
	return encryptedPrefix + base64.RawURLEncoding.EncodeToString(sealed), nil
}

// Open decrypts a value produced by Seal for the same context. Errors never
// include the secret.
func (b *SecretBox) Open(stored, context string) (string, error) {
	if !IsEncrypted(stored) {
		return "", errors.New("value is not encrypted")
	}
	data, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(stored, encryptedPrefix))
	if err != nil {
		return "", errors.New("malformed encrypted secret")
	}
	if len(data) < b.aead.NonceSize()+b.aead.Overhead() {
		return "", errors.New("malformed encrypted secret")
	}
	nonce, ciphertext := data[:b.aead.NonceSize()], data[b.aead.NonceSize():]
	plaintext, err := b.aead.Open(nil, nonce, ciphertext, []byte(context))
	if err != nil {
		return "", errors.New("cannot decrypt secret: wrong key or corrupted value")
	}
	return string(plaintext), nil
}
