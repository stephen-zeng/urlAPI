package cmd

import (
	"path/filepath"
	"testing"
	"time"
	"urlAPI/internal/auth"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

func TestResetPasswordGeneratesUniqueCredential(t *testing.T) {
	gdb, err := database.Open(filepath.Join(t.TempDir(), "admin.db"))
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
	})
	database.SetLocalDB(gdb)
	s := util.DefaultAppSettings()
	s.Security.DashboardPasswordHash = auth.ClientCredential("123456")
	if err := database.SaveAppSettings(s); err != nil {
		t.Fatal(err)
	}
	session := model.Session{Token: "existing", Expire: time.Now().Add(time.Hour)}
	if err := database.CreateSession(&session); err != nil {
		t.Fatal(err)
	}

	password, err := resetPassword()
	if err != nil {
		t.Fatal(err)
	}
	if password == "" || password == "123456" {
		t.Fatalf("reset produced %q", password)
	}
	hash := database.SettingsStore.Get().Security.DashboardPasswordHash
	if ok, _ := auth.VerifyPassword(hash, auth.ClientCredential(password)); !ok {
		t.Fatal("new password does not verify")
	}
	if ok, _ := auth.VerifyPassword(hash, auth.ClientCredential("123456")); ok {
		t.Fatal("old default password still works")
	}
	if _, ok := database.Sessions.Get("existing"); ok {
		t.Fatal("sessions not cleared")
	}
	second, err := resetPassword()
	if err != nil || second == password {
		t.Fatalf("second reset returned %q, %v", second, err)
	}
}
