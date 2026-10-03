package op

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"urlAPI/internal/auth"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

// setupSettingsDB opens a fully migrated database with the given settings.
func setupSettingsDB(t *testing.T, settings util.AppSettings) {
	t.Helper()
	gdb, err := database.Open(filepath.Join(t.TempDir(), "auth.db"))
	if err != nil {
		t.Fatal(err)
	}
	prevDB, prevSettings := database.GetLocalDB().DB(), database.SettingsStore.Get()
	t.Cleanup(func() {
		if sqlDB, err := gdb.DB(); err == nil {
			_ = sqlDB.Close()
		}
		database.SetLocalDB(prevDB)
		database.SettingsStore.Replace(prevSettings)
		database.SetSecretBox(nil)
		db = nil
	})
	db = database.SetLocalDB(gdb)
	if err := database.SaveAppSettings(settings); err != nil {
		t.Fatal(err)
	}
}

func loginWith(credential string) (Session, error) {
	return HandleSession(Session{Operation: "login"}, model.Session{Token: credential})
}

func TestLoginWithModernHash(t *testing.T) {
	s := util.DefaultAppSettings()
	hash, err := auth.HashPassword(auth.ClientCredential("hunter2"))
	if err != nil {
		t.Fatal(err)
	}
	s.Security.DashboardPasswordHash = hash
	setupSettingsDB(t, s)

	resp, err := loginWith(auth.ClientCredential("hunter2"))
	if err != nil || len(resp.SessionToken) != 64 {
		t.Fatalf("login: %+v %v", resp, err)
	}
	if _, ok := database.Sessions.Get(resp.SessionToken); !ok {
		t.Fatal("session not stored")
	}
	if _, err := loginWith(auth.ClientCredential("hunter3")); err == nil {
		t.Fatal("wrong password accepted")
	}
	if _, err := loginWith(hash); err == nil {
		t.Fatal("the stored hash itself must not work as a credential")
	}
	// Re-validating with the session token (as the dashboard does on load)
	// succeeds without issuing a new session.
	before := database.Sessions.Len()
	if _, err := loginWith(resp.SessionToken); err != nil {
		t.Fatalf("session revalidation: %v", err)
	}
	if database.Sessions.Len() != before {
		t.Fatal("revalidation issued a new session")
	}
}

func TestLegacyPasswordMigratesOnLogin(t *testing.T) {
	legacy := auth.ClientCredential("legacy-pass")
	s := util.DefaultAppSettings()
	s.Security.DashboardPasswordHash = legacy
	setupSettingsDB(t, s)

	if _, err := loginWith(auth.ClientCredential("nope")); err == nil {
		t.Fatal("wrong password accepted")
	}
	if got := database.SettingsStore.Get().Security.DashboardPasswordHash; got != legacy {
		t.Fatal("failed login must not touch the stored hash")
	}

	if _, err := loginWith(legacy); err != nil {
		t.Fatalf("legacy login failed: %v", err)
	}
	migrated := database.SettingsStore.Get().Security.DashboardPasswordHash
	if !strings.HasPrefix(migrated, "$argon2id$") {
		t.Fatalf("hash not migrated: %q", migrated)
	}
	// Persisted, not only in memory.
	var row model.AppSetting
	if err := database.GetLocalDB().DB().Where("key = ?", "security.dashboard_password_hash").First(&row).Error; err != nil || row.Value != migrated {
		t.Fatalf("migrated hash not persisted: %q %v", row.Value, err)
	}
	if _, err := loginWith(legacy); err != nil {
		t.Fatalf("login after migration failed: %v", err)
	}
}

func TestPasswordChangeStoresArgon2(t *testing.T) {
	s := util.DefaultAppSettings()
	s.Security.DashboardPasswordHash = auth.ClientCredential("old")
	setupSettingsDB(t, s)

	body, _ := json.Marshal(map[string]any{
		"password_hash":         auth.ClientCredential("new-password"),
		"dashboard_allowed_ips": []string{"*"},
		"allowed_referers":      []string{"localhost"},
	})
	info := Session{SettingPart: "security", SettingBody: body}
	if err := editSettings(&info); err != nil {
		t.Fatal(err)
	}
	stored := database.SettingsStore.Get().Security.DashboardPasswordHash
	if !strings.HasPrefix(stored, "$argon2id$") || strings.Contains(stored, auth.ClientCredential("new-password")) {
		t.Fatalf("password stored as %q", stored)
	}
	if ok, _ := auth.VerifyPassword(stored, auth.ClientCredential("new-password")); !ok {
		t.Fatal("new password does not verify")
	}
}

func TestSettingsAPIDoesNotExposeSecrets(t *testing.T) {
	s := util.DefaultAppSettings()
	hash, _ := auth.HashPassword(auth.ClientCredential("pw"))
	s.Security.DashboardPasswordHash = hash
	s.Providers.OpenAI.APIKey = "sk-should-never-leak"
	s.Providers.Alibaba.APIKey = "ali-should-never-leak"
	s.Web.RepoToken = "ghp-should-never-leak"
	s.Web.YouTubeToken = "yt-should-never-leak"
	setupSettingsDB(t, s)

	for _, part := range []string{"openai", "alibaba", "deepseek", "otherapi", "web", "security", "txt", "img", "rand", "contxt", "taskBehavior", "txtSecurity", "imgSecurity"} {
		info := Session{SettingPart: part}
		if err := fetchSettings(&info); err != nil {
			t.Fatalf("%s: %v", part, err)
		}
		body := string(info.SettingBody)
		for _, secret := range []string{"should-never-leak", hash, "$argon2id$"} {
			if strings.Contains(body, secret) {
				t.Fatalf("settings part %s exposes a secret: %s", part, body)
			}
		}
	}
	info := Session{SettingPart: "openai"}
	_ = fetchSettings(&info)
	if !strings.Contains(string(info.SettingBody), `"api_key_set":true`) {
		t.Fatalf("api_key_set missing: %s", info.SettingBody)
	}
}

func TestExpiredSessionIsRejected(t *testing.T) {
	setupSettingsDB(t, util.DefaultAppSettings())
	database.Sessions.Set(model.Session{Token: "expired", Expire: time.Now().Add(-time.Minute)})
	t.Cleanup(func() { database.Sessions.Delete("expired") })
	if _, err := HandleSession(Session{Operation: "fetchRepo"}, model.Session{Token: "expired"}); err == nil {
		t.Fatal("expired session accepted")
	}
}
