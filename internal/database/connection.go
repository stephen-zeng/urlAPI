package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/common-nighthawk/go-figure"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SQLite settings for a concurrent HTTP server:
//
//   - WAL lets readers proceed while a single writer commits, instead of the
//     rollback journal's exclusive database lock. The setting is persistent
//     in the database file and remains readable by older builds. WAL needs a
//     local filesystem (shared memory); do not put the database on NFS/SMB.
//   - synchronous=NORMAL is the recommended pairing with WAL: the database
//     cannot be corrupted by a crash, but the most recent transactions may
//     be lost on power failure. Task logs, sessions and caches tolerate that,
//     and it avoids an fsync per logged request.
//   - busy_timeout makes a connection wait for the write lock instead of
//     failing immediately with SQLITE_BUSY, since SQLite allows only one
//     writer at a time.
//   - _txlock=immediate takes the write lock when a transaction starts, so a
//     read-then-write transaction waits (honouring busy_timeout) rather than
//     failing with a lock-upgrade deadlock.
const (
	sqliteBusyTimeoutMS = 5000
	sqliteMaxOpenConns  = 8
)

func sqliteDSN(path string) string {
	return fmt.Sprintf("%s?_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=%d&_txlock=immediate", path, sqliteBusyTimeoutMS)
}

// Open opens (creating if needed) the SQLite database at path, applies the
// connection settings above and migrates the schema. On error nothing is left
// open.
func Open(path string) (*gorm.DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, errors.Wrap(err, "create database directory")
		}
	}
	db, err := gorm.Open(sqlite.Open(sqliteDSN(path)), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "gorm")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "sql db")
	}
	sqlDB.SetMaxOpenConns(sqliteMaxOpenConns)
	sqlDB.SetMaxIdleConns(sqliteMaxOpenConns)
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, errors.Wrap(err, "ping database")
	}
	if err := migrate(db); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return db, nil
}

// 包括所有数据的初始化
func Init() error {
	figlet := figure.NewFigure("urlAPI", "", true)
	figlet.Print()
	db, err := Open(dbPath)
	if err != nil {
		return err
	}
	SetLocalDB(db)
	log.Println("Connected to database")
	if err := loadState(); err != nil {
		_ = Disconnect()
		return err
	}
	return nil
}

func loadState() error {
	if err := initRepoMap(); err != nil {
		return err
	}
	if err := initSessionMap(); err != nil {
		return err
	}
	if err := initAppSettings(); err != nil {
		return errors.Wrap(err, "initAppSettings")
	}
	return nil
}

func migrate(db *gorm.DB) error {
	models := []any{&AppSetting{}, &Provider{}, &ServiceConfig{}, &Prompt{}, &ConfigListItem{}, &Task{}, &Session{}, &Repo{}}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return errors.Wrapf(err, "migrate %T", model)
		}
	}
	return nil
}

// Disconnect closes the database. It is safe to call when no database is
// open and more than once.
func Disconnect() error {
	if localDB.db == nil {
		return nil
	}
	sqlDB, err := localDB.db.DB()
	localDB.db = nil
	if err != nil {
		return errors.Wrap(err, "sql db")
	}
	if err := sqlDB.Close(); err != nil {
		return errors.Wrap(err, "close database")
	}
	log.Println("Disconnected from database")
	return nil
}

func initRepoMap() error {
	var repos []Repo
	if err := localDB.db.Find(&repos).Error; err != nil {
		return errors.Wrap(err, "db find")
	}
	for _, repo := range repos {
		var repoList []string
		if err := json.Unmarshal([]byte(repo.Content), &repoList); err != nil {
			return errors.Wrap(err, "json")
		}
		Repos.Set(repo.API, repo.Info, repoList)
	}
	log.Println("Initialized RepoMap")
	return nil
}

func initSessionMap() error {
	var sessions []Session
	if err := localDB.db.Find(&sessions).Error; err != nil {
		return errors.Wrap(err, "db")
	}
	Sessions.Reset(sessions)
	log.Println("Initialized SessionMap")
	return nil
}

func ClearTask() {
	if localDB.db.Migrator().HasTable(&Task{}) {
		if err := localDB.db.Migrator().DropTable(&Task{}); err != nil {
			log.Fatal(errors.Wrap(err, "db"))
		}
		if err := localDB.db.AutoMigrate(&Task{}); err != nil {
			log.Fatal(errors.Wrap(err, "db"))
		}
	}
}

func ClearSession() {
	if localDB.db.Migrator().HasTable(&Session{}) {
		if err := localDB.db.Migrator().DropTable(&Session{}); err != nil {
			log.Fatal(errors.Wrap(err, "db"))
		}
		if err := localDB.db.AutoMigrate(&Session{}); err != nil {
			log.Fatal(errors.Wrap(err, "db"))
		}
	}
	Sessions.Reset(nil)
}
