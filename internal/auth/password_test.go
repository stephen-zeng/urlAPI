package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

func TestHashAndVerifyPassword(t *testing.T) {
	cred := ClientCredential("correct horse battery staple")
	hash, err := HashPassword(cred)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(hash, "$argon2id$v=19$") || strings.Contains(hash, cred) {
		t.Fatalf("unexpected hash format: %s", hash)
	}
	if ok, rehash := VerifyPassword(hash, cred); !ok || rehash {
		t.Fatalf("VerifyPassword(correct) = %v, %v", ok, rehash)
	}
	if ok, _ := VerifyPassword(hash, ClientCredential("wrong")); ok {
		t.Fatal("wrong password accepted")
	}
	if ok, _ := VerifyPassword(hash, ""); ok {
		t.Fatal("empty credential accepted")
	}
	other, err := HashPassword(cred)
	if err != nil {
		t.Fatal(err)
	}
	if other == hash {
		t.Fatal("hashes are not salted")
	}
	if _, err := HashPassword(""); err == nil {
		t.Fatal("empty password must be rejected")
	}
}

func TestLegacySHA256Verification(t *testing.T) {
	sum := sha256.Sum256([]byte("123456"))
	legacy := hex.EncodeToString(sum[:])
	if legacy != "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92" {
		t.Fatal("fixture mismatch")
	}
	if !IsLegacyHash(legacy) || !IsDefaultCredential(legacy) {
		t.Fatal("legacy default hash not recognised")
	}
	if ok, rehash := VerifyPassword(legacy, ClientCredential("123456")); !ok || !rehash {
		t.Fatalf("legacy verify = %v, %v; want true, true", ok, rehash)
	}
	if ok, _ := VerifyPassword(legacy, ClientCredential("1234567")); ok {
		t.Fatal("wrong password accepted for legacy hash")
	}
	if ok, _ := VerifyPassword(legacy, strings.ToUpper(legacy)); !ok {
		t.Fatal("hex case should not matter for legacy credentials")
	}
}

func TestVerifyRejectsMalformedStoredValues(t *testing.T) {
	cred := ClientCredential("pw")
	for _, stored := range []string{
		"", "plaintext", "$argon2id$", "$argon2i$v=19$m=19456,t=2,p=1$c2FsdHNhbHQ$a2V5a2V5a2V5a2V5a2V5",
		"$argon2id$v=18$m=19456,t=2,p=1$c2FsdHNhbHQ$a2V5a2V5a2V5a2V5a2V5",
		"$argon2id$v=19$m=999999999,t=2,p=1$c2FsdHNhbHQ$a2V5a2V5a2V5a2V5a2V5",
		"$argon2id$v=19$m=19456,t=2,p=1$!!$a2V5",
	} {
		if ok, _ := VerifyPassword(stored, cred); ok {
			t.Errorf("VerifyPassword(%q) accepted", stored)
		}
	}
}

func TestOutdatedParametersNeedRehash(t *testing.T) {
	cred := ClientCredential("pw")
	salt := []byte("0123456789abcdef")
	key := deriveKey(cred, salt, 1, 8*1024, 1, 32)
	stored := "$argon2id$v=19$m=8192,t=1,p=1$" + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(key)
	if ok, rehash := VerifyPassword(stored, cred); !ok || !rehash {
		t.Fatalf("outdated params: %v, %v", ok, rehash)
	}
}

func TestGeneratePassword(t *testing.T) {
	pw, err := GeneratePassword()
	if err != nil {
		t.Fatal(err)
	}
	if len(pw) != 24 || strings.ContainsAny(pw, "+/= ") || pw == "123456" {
		t.Fatalf("unexpected generated password %q", pw)
	}
}
