// Package auth implements administrator password hashing and the
// authenticated encryption of stored provider secrets.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

// The dashboard never sends the raw password: it sends the hex SHA-256 of
// it (the "client credential"). Earlier versions stored that value as is,
// which made the database row itself a usable credential. The server now
// stores an Argon2id hash of the client credential instead.

// EnvAdminPassword optionally sets the initial dashboard password of a new
// installation. Without it a random password is generated and printed once.
const EnvAdminPassword = "URLAPI_ADMIN_PASSWORD"

// Argon2id parameters (OWASP baseline: 19 MiB, 2 iterations, 1 lane).
const (
	argonMemoryKiB = 19 * 1024
	argonTime      = 2
	argonThreads   = 1
	argonKeyLen    = 32
	argonSaltLen   = 16
)

var legacyHashPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

// knownDefaultCredential is the client credential for the historical
// default password "123456".
var knownDefaultCredential = ClientCredential("123456")

// verifySlots bounds concurrent Argon2id computations so a burst of login
// attempts cannot exhaust memory.
var verifySlots = make(chan struct{}, 4)

// ClientCredential returns what the dashboard sends for password: the
// lower-case hex SHA-256 digest.
func ClientCredential(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

// IsDefaultCredential reports whether credential is the historical default
// password.
func IsDefaultCredential(credential string) bool {
	return subtle.ConstantTimeCompare([]byte(credential), []byte(knownDefaultCredential)) == 1
}

// IsLegacyHash reports whether stored is an unsalted SHA-256 hash from an
// earlier version.
func IsLegacyHash(stored string) bool {
	return legacyHashPattern.MatchString(stored)
}

// HashPassword returns an encoded Argon2id hash of a client credential.
func HashPassword(credential string) (string, error) {
	if credential == "" {
		return "", errors.New("empty password")
	}
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", errors.Wrap(err, "generate salt")
	}
	key := deriveKey(credential, salt, argonTime, argonMemoryKiB, argonThreads, argonKeyLen)
	b64 := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemoryKiB, argonTime, argonThreads, b64.EncodeToString(salt), b64.EncodeToString(key)), nil
}

func deriveKey(credential string, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte {
	verifySlots <- struct{}{}
	defer func() { <-verifySlots }()
	return argon2.IDKey([]byte(credential), salt, time, memory, threads, keyLen)
}

// VerifyPassword checks a client credential against a stored hash. It
// accepts Argon2id hashes and legacy unsalted SHA-256 hashes; needsRehash
// is true when the stored hash should be replaced with HashPassword's
// output. An empty or unrecognised stored value never matches.
func VerifyPassword(stored, credential string) (ok, needsRehash bool) {
	if stored == "" || credential == "" {
		return false, false
	}
	if IsLegacyHash(stored) {
		match := subtle.ConstantTimeCompare([]byte(stored), []byte(strings.ToLower(credential))) == 1
		return match, match
	}
	params, salt, want, err := decodeArgon2id(stored)
	if err != nil {
		return false, false
	}
	got := deriveKey(credential, salt, params.time, params.memory, params.threads, uint32(len(want)))
	if subtle.ConstantTimeCompare(got, want) != 1 {
		return false, false
	}
	current := params.time == argonTime && params.memory == argonMemoryKiB && params.threads == argonThreads && len(want) == argonKeyLen
	return true, !current
}

type argonParams struct {
	time    uint32
	memory  uint32
	threads uint8
}

func decodeArgon2id(encoded string) (argonParams, []byte, []byte, error) {
	var p argonParams
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return p, nil, nil, errors.New("not an argon2id hash")
	}
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return p, nil, nil, errors.New("unsupported argon2 version")
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.time, &p.threads); err != nil {
		return p, nil, nil, errors.New("invalid argon2 parameters")
	}
	// Refuse absurd parameters from a tampered database.
	if p.time < 1 || p.time > 10 || p.memory < 8*1024 || p.memory > 1024*1024 || p.threads < 1 || p.threads > 16 {
		return p, nil, nil, errors.New("argon2 parameters out of range")
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) < 8 {
		return p, nil, nil, errors.New("invalid argon2 salt")
	}
	key, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(key) < 16 || len(key) > 64 {
		return p, nil, nil, errors.New("invalid argon2 key")
	}
	return p, salt, key, nil
}

// GeneratePassword returns a random password with 144 bits of entropy.
func GeneratePassword() (string, error) {
	buf := make([]byte, 18)
	if _, err := rand.Read(buf); err != nil {
		return "", errors.Wrap(err, "generate password")
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
