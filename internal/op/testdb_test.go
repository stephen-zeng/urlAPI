package op

import (
	"path/filepath"
	"testing"
	"urlAPI/internal/database"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newTestDB wires the op package to a fresh temporary SQLite database.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "test.db")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := gdb.AutoMigrate(&database.Task{}, &database.Session{}, &database.Repo{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	prev := database.GetLocalDB().DB()
	t.Cleanup(func() {
		if sqlDB, err := gdb.DB(); err == nil {
			_ = sqlDB.Close()
		}
		database.SetLocalDB(prev)
		db = nil
	})
	db = database.SetLocalDB(gdb)
	return gdb
}
