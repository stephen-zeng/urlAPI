package database

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"urlAPI/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func pragma(t *testing.T, db *gorm.DB, name string) string {
	t.Helper()
	var value string
	if err := db.Raw("PRAGMA " + name).Scan(&value).Error; err != nil {
		t.Fatalf("PRAGMA %s: %v", name, err)
	}
	return value
}

func closeDB(t *testing.T, db *gorm.DB) {
	t.Helper()
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestOpenConfiguresSQLite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "dir", "database.db")
	db, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer closeDB(t, db)

	info, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm&0o002 != 0 {
		t.Fatalf("database directory is world-writable: %v", perm)
	}
	if got := pragma(t, db, "journal_mode"); got != "wal" {
		t.Fatalf("journal_mode = %q, want wal", got)
	}
	if got := pragma(t, db, "busy_timeout"); got != "5000" {
		t.Fatalf("busy_timeout = %q, want 5000", got)
	}
	if got := pragma(t, db, "synchronous"); got != "1" { // NORMAL
		t.Fatalf("synchronous = %q, want 1 (NORMAL)", got)
	}
	for _, table := range []string{"tasks", "sessions", "repos", "app_settings", "providers", "service_configs", "prompts", "config_list_items"} {
		if !db.Migrator().HasTable(table) {
			t.Errorf("table %s not migrated", table)
		}
	}
	if !db.Migrator().HasIndex(&model.Task{}, "Time") {
		t.Error("tasks.time index missing")
	}
}

func TestOpenExistingLegacyDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.db")
	// Simulate a database created by an older build: default journal mode,
	// no tasks.time index, existing rows.
	legacy, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	stmts := []string{
		"CREATE TABLE tasks (uuid text PRIMARY KEY, time datetime, ip text, type text, status text, target text, `return` text, region text, referer text, device text, more_info text, api text, model text, temp text, size text)",
		"CREATE TABLE sessions (token text PRIMARY KEY, expire datetime, term numeric)",
		"INSERT INTO tasks (uuid, time, type) VALUES ('old', '2024-01-02 03:04:05+00:00', 'txt')",
		"INSERT INTO sessions (token, expire, term) VALUES ('legacy-token', '2999-01-01 00:00:00+00:00', 1)",
	}
	for _, stmt := range stmts {
		if err := legacy.Exec(stmt).Error; err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	closeDB(t, legacy)

	for i := 0; i < 2; i++ { // reopening must be idempotent
		db, err := Open(path)
		if err != nil {
			t.Fatalf("Open #%d: %v", i, err)
		}
		var task model.Task
		if err := db.First(&task, "uuid = ?", "old").Error; err != nil || task.Type != "txt" {
			t.Fatalf("legacy task lost: %+v %v", task, err)
		}
		var session model.Session
		if err := db.First(&session, "token = ?", "legacy-token").Error; err != nil || !session.Expire.After(time.Now()) {
			t.Fatalf("legacy session lost: %+v %v", session, err)
		}
		closeDB(t, db)
	}
}

func TestOpenFailsWhenDirectoryCannotBeCreated(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocker, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(filepath.Join(blocker, "db", "database.db")); err == nil {
		t.Fatal("expected error when the parent path is a file")
	}
}

func TestInitAndDisconnectLifecycle(t *testing.T) {
	origPath, origDB := dbPath, localDB.db
	t.Cleanup(func() { dbPath, localDB.db = origPath, origDB })
	dbPath = filepath.Join(t.TempDir(), "assets", "database.db")
	localDB.db = nil

	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if localDB.db == nil {
		t.Fatal("Init did not publish the database")
	}
	if err := Disconnect(); err != nil {
		t.Fatalf("Disconnect: %v", err)
	}
	if err := Disconnect(); err != nil {
		t.Fatalf("second Disconnect: %v", err)
	}
	if localDB.db != nil {
		t.Fatal("Disconnect left a stale handle")
	}
}

func TestInitFailureLeavesNoGlobalState(t *testing.T) {
	origPath, origDB := dbPath, localDB.db
	t.Cleanup(func() { dbPath, localDB.db = origPath, origDB })
	dir := t.TempDir()
	dbPath = filepath.Join(dir, "database.db")
	localDB.db = nil

	// A corrupt repo row makes loading state fail after the open succeeds.
	db, err := Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO repos (uuid, api, info, content) VALUES ('r', 'github', 'a/b', 'not json')").Error; err != nil {
		t.Fatal(err)
	}
	closeDB(t, db)

	err = Init()
	if err == nil || !strings.Contains(err.Error(), "json") {
		t.Fatalf("expected json error, got %v", err)
	}
	if localDB.db != nil {
		t.Fatal("failed Init left a database handle published")
	}
}

func TestConcurrentWritesDoNotFailBusy(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "busy.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer closeDB(t, db)
	adapter := &SQLiteAdapter{db: db}
	const writers = 16
	errs := make(chan error, writers)
	for w := 0; w < writers; w++ {
		go func(w int) {
			var err error
			for i := 0; i < 20 && err == nil; i++ {
				err = adapter.CreateTask(&model.Task{UUID: time.Now().Format(time.RFC3339Nano) + string(rune('a'+w)) + string(rune('A'+i)), Time: time.Now()})
			}
			errs <- err
		}(w)
	}
	for w := 0; w < writers; w++ {
		if err := <-errs; err != nil {
			t.Fatalf("concurrent write failed: %v", err)
		}
	}
	var count int64
	db.Model(&model.Task{}).Count(&count)
	if count != writers*20 {
		t.Fatalf("count = %d, want %d", count, writers*20)
	}
}
