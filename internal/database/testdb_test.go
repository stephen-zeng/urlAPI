package database

import (
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newTestDB points the package adapter at a fresh SQLite database in a
// temporary directory and restores the previous adapter afterwards.
func newTestDB(t *testing.T) *SQLiteAdapter {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "test.db")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := gdb.AutoMigrate(&Task{}, &Session{}, &Repo{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	prev := localDB.db
	t.Cleanup(func() {
		if sqlDB, err := gdb.DB(); err == nil {
			_ = sqlDB.Close()
		}
		localDB.db = prev
	})
	return SetLocalDB(gdb)
}
