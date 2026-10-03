package database

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"urlAPI/internal/auth"
	"urlAPI/internal/model"
	"urlAPI/util"
)

const (
	plainKey   = "sk-provider-secret-0001"
	plainToken = "ghp_repoTokenSecret0001"
)

// openSettingsDB opens a fully migrated temporary database as the package
// database and isolates the secret codec and settings store.
func openSettingsDB(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "settings.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	prevDB, prevSettings := localDB.db, SettingsStore.Get()
	secrets.mu.Lock()
	prevBox, prevPreserved := secrets.box, secrets.preserved
	secrets.box, secrets.preserved = nil, make(map[string]string)
	secrets.mu.Unlock()
	t.Cleanup(func() {
		closeDB(t, db)
		localDB.db = prevDB
		SettingsStore.Replace(prevSettings)
		secrets.mu.Lock()
		secrets.box, secrets.preserved = prevBox, prevPreserved
		secrets.mu.Unlock()
	})
	SetLocalDB(db)
	return path
}

func box(t *testing.T, fill byte) *auth.SecretBox {
	t.Helper()
	b, err := auth.NewSecretBox(bytes.Repeat([]byte{fill}, 32))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func storedProviderKey(t *testing.T, name string) string {
	t.Helper()
	var p model.Provider
	if err := localDB.db.Where("name = ?", name).First(&p).Error; err != nil {
		t.Fatal(err)
	}
	return p.APIKeyEnc
}

func storedWebSettings(t *testing.T) map[string]string {
	t.Helper()
	var s model.ServiceConfig
	if err := localDB.db.Where("service = ?", "web").First(&s).Error; err != nil {
		t.Fatal(err)
	}
	values := map[string]string{}
	if err := json.Unmarshal([]byte(s.Settings), &values); err != nil {
		t.Fatal(err)
	}
	return values
}

func settingsWithSecrets() util.AppSettings {
	s := util.DefaultAppSettings()
	s.Security.DashboardPasswordHash = "$argon2id$placeholder"
	s.Providers.OpenAI.APIKey = plainKey
	s.Web.RepoToken = plainToken
	return s
}

func reload(t *testing.T) util.AppSettings {
	t.Helper()
	s, err := loadAppSettings()
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSecretsWithoutKeyUseLegacyFormat(t *testing.T) {
	openSettingsDB(t)
	if err := SaveAppSettings(settingsWithSecrets()); err != nil {
		t.Fatal(err)
	}
	if got := storedProviderKey(t, "openai"); got != base64.StdEncoding.EncodeToString([]byte(plainKey)) {
		t.Fatalf("stored provider key = %q", got)
	}
	if got := storedWebSettings(t)["repo_token"]; got != plainToken {
		t.Fatalf("stored repo token = %q", got)
	}
	if s := reload(t); s.Providers.OpenAI.APIKey != plainKey || s.Web.RepoToken != plainToken {
		t.Fatalf("legacy values not readable: %+v", s.Providers.OpenAI)
	}
}

func TestLegacyBase64SecretsMigrateWhenKeyConfigured(t *testing.T) {
	openSettingsDB(t)
	// Write as an older version would (no key: Base64 / plain text).
	if err := SaveAppSettings(settingsWithSecrets()); err != nil {
		t.Fatal(err)
	}
	// Also a row from before Base64 was introduced (raw value).
	if err := localDB.db.Model(&model.Provider{}).Where("name = ?", "deepseek").Update("api_key_enc", "raw-legacy-key!").Error; err != nil {
		t.Fatal(err)
	}

	SetSecretBox(box(t, 9))
	if err := initAppSettings(); err != nil { // startup re-save performs the migration
		t.Fatal(err)
	}
	for _, name := range []string{"openai", "deepseek"} {
		if got := storedProviderKey(t, name); !auth.IsEncrypted(got) {
			t.Fatalf("%s key not encrypted after migration: %q", name, got)
		}
	}
	web := storedWebSettings(t)
	if !auth.IsEncrypted(web["repo_token"]) || strings.Contains(web["repo_token"], plainToken) {
		t.Fatalf("repo token not encrypted: %q", web["repo_token"])
	}
	if web["youtube_token"] != "" {
		t.Fatalf("empty token should stay empty, got %q", web["youtube_token"])
	}
	s := reload(t)
	if s.Providers.OpenAI.APIKey != plainKey || s.Providers.DeepSeek.APIKey != "raw-legacy-key!" || s.Web.RepoToken != plainToken {
		t.Fatalf("decrypted values wrong: %q %q %q", s.Providers.OpenAI.APIKey, s.Providers.DeepSeek.APIKey, s.Web.RepoToken)
	}
}

func TestEncryptedSecretsSurviveMissingOrWrongKey(t *testing.T) {
	openSettingsDB(t)
	SetSecretBox(box(t, 9))
	if err := SaveAppSettings(settingsWithSecrets()); err != nil {
		t.Fatal(err)
	}
	original := storedProviderKey(t, "openai")
	originalToken := storedWebSettings(t)["repo_token"]

	for name, b := range map[string]*auth.SecretBox{"missing key": nil, "wrong key": box(t, 8)} {
		t.Run(name, func(t *testing.T) {
			SetSecretBox(b)
			s := reload(t)
			if s.Providers.OpenAI.APIKey != "" || s.Web.RepoToken != "" {
				t.Fatal("undecryptable secret must not be exposed")
			}
			// Saving (as every startup does) must not destroy the ciphertext.
			if err := SaveAppSettings(s); err != nil {
				t.Fatal(err)
			}
			if got := storedProviderKey(t, "openai"); got != original {
				t.Fatalf("ciphertext changed: %q -> %q", original, got)
			}
			if got := storedWebSettings(t)["repo_token"]; got != originalToken {
				t.Fatalf("token ciphertext changed")
			}
		})
	}

	// With the right key again everything is readable.
	SetSecretBox(box(t, 9))
	if s := reload(t); s.Providers.OpenAI.APIKey != plainKey || s.Web.RepoToken != plainToken {
		t.Fatal("secrets lost after key restored")
	}
}

func TestNewSecretReplacesUndecryptableValue(t *testing.T) {
	openSettingsDB(t)
	SetSecretBox(box(t, 9))
	if err := SaveAppSettings(settingsWithSecrets()); err != nil {
		t.Fatal(err)
	}
	SetSecretBox(box(t, 8)) // key rotated without migration
	s := reload(t)
	s.Providers.OpenAI.APIKey = "sk-new-key"
	if err := SaveAppSettings(s); err != nil {
		t.Fatal(err)
	}
	if got := reload(t).Providers.OpenAI.APIKey; got != "sk-new-key" {
		t.Fatalf("new key not stored: %q", got)
	}
}

func TestInitBootstrapsUniqueAdminPassword(t *testing.T) {
	openSettingsDB(t)
	t.Setenv(auth.EnvAdminPassword, "")
	if err := initAppSettings(); err != nil {
		t.Fatal(err)
	}
	hash := SettingsStore.Get().Security.DashboardPasswordHash
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatalf("new installation password hash = %q", hash)
	}
	if ok, _ := auth.VerifyPassword(hash, auth.ClientCredential("123456")); ok {
		t.Fatal("new installation uses the well-known default password")
	}
	// A second start keeps the same password.
	if err := initAppSettings(); err != nil {
		t.Fatal(err)
	}
	if SettingsStore.Get().Security.DashboardPasswordHash != hash {
		t.Fatal("password changed on restart")
	}
}

func TestInitUsesConfiguredAdminPassword(t *testing.T) {
	openSettingsDB(t)
	t.Setenv(auth.EnvAdminPassword, "s3cret-initial")
	if err := initAppSettings(); err != nil {
		t.Fatal(err)
	}
	if ok, _ := auth.VerifyPassword(SettingsStore.Get().Security.DashboardPasswordHash, auth.ClientCredential("s3cret-initial")); !ok {
		t.Fatal("configured initial password not applied")
	}
}

func TestInitKeepsExistingLegacyPassword(t *testing.T) {
	openSettingsDB(t)
	legacy := auth.ClientCredential("my-old-password")
	s := util.DefaultAppSettings()
	s.Security.DashboardPasswordHash = legacy
	if err := SaveAppSettings(s); err != nil {
		t.Fatal(err)
	}
	t.Setenv(auth.EnvAdminPassword, "ignored-for-existing-installs")
	if err := initAppSettings(); err != nil {
		t.Fatal(err)
	}
	if got := SettingsStore.Get().Security.DashboardPasswordHash; got != legacy {
		t.Fatalf("existing password replaced: %q", got)
	}
}
