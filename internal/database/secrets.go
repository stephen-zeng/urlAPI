package database

import (
	"encoding/base64"
	"log"
	"sync"
	"urlAPI/internal/auth"
)

// Stored secret fields and their encryption contexts.
const (
	webRepoTokenKey    = "repo_token"
	webYouTubeTokenKey = "youtube_token"
)

func providerSecretContext(name string) string { return "provider/" + name + "/api_key" }
func webSecretContext(key string) string       { return "service/web/" + key }

// secretCodec encrypts secrets on write when a master key is configured and
// decrypts them on read. Without a key, values are written in the legacy
// format (Base64 for provider keys, plain text for web tokens) so existing
// installations keep working; configuring a key migrates them on the next
// save, which happens at every startup.
//
// Encrypted values that cannot be decrypted (missing or wrong key) are kept
// verbatim and written back unchanged, so a misconfiguration never destroys
// stored credentials.
type secretCodec struct {
	mu        sync.Mutex
	box       *auth.SecretBox
	preserved map[string]string
}

var secrets = &secretCodec{preserved: make(map[string]string)}

// SetSecretBox installs the master key used for stored secrets (nil for
// none).
func SetSecretBox(box *auth.SecretBox) {
	secrets.mu.Lock()
	defer secrets.mu.Unlock()
	secrets.box = box
}

func (c *secretCodec) encrypting() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.box != nil
}

// open returns the plaintext of a stored value, or "" if it cannot be
// decrypted, in which case the stored value is preserved for writing back.
func (c *secretCodec) open(stored, context string, legacy func(string) string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.preserved, context)
	if stored == "" {
		return ""
	}
	if !auth.IsEncrypted(stored) {
		return legacy(stored)
	}
	if c.box == nil {
		log.Printf("Secret %s is encrypted but %s is not set; it is unavailable until the key is configured", context, auth.EnvSecretKey)
		c.preserved[context] = stored
		return ""
	}
	plaintext, err := c.box.Open(stored, context)
	if err != nil {
		log.Printf("Secret %s: %v; keeping the stored value unchanged", context, err)
		c.preserved[context] = stored
		return ""
	}
	return plaintext
}

// seal returns the value to store for plaintext.
func (c *secretCodec) seal(plaintext, context string, legacy func(string) string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if plaintext == "" {
		// Nothing new was provided: keep an undecryptable original.
		return c.preserved[context], nil
	}
	delete(c.preserved, context)
	if c.box == nil {
		return legacy(plaintext), nil
	}
	return c.box.Seal(plaintext, context)
}

func identity(s string) string { return s }

func encodeLegacyBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// decodeLegacyBase64 reads the historical Base64 encoding, falling back to
// the raw value for rows written before it was introduced.
func decodeLegacyBase64(value string) string {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return value
	}
	return string(decoded)
}
