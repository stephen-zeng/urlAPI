package database

import (
	"database/sql"
	"encoding/json"
	"github.com/common-nighthawk/go-figure"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

// 包括所有数据的初始化
func Init() error {
	figlet := figure.NewFigure("urlAPI", "", true)
	figlet.Print()
	if err := connect(); err != nil {
		return err
	}
	migration()
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

func migration() {
	localDB.db.AutoMigrate(&AppSetting{})
	localDB.db.AutoMigrate(&Provider{})
	localDB.db.AutoMigrate(&ServiceConfig{})
	localDB.db.AutoMigrate(&Prompt{})
	localDB.db.AutoMigrate(&ConfigListItem{})
	localDB.db.AutoMigrate(&Task{})
	localDB.db.AutoMigrate(&Session{})
	localDB.db.AutoMigrate(&Repo{})
}

func connect() error {
	var err error
	os.Mkdir("assets", 0777)
	tmp, _ := sql.Open("sqlite3", dbPath)
	tmp.Close()
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return errors.Wrap(err, "gorm")
	}
	SetLocalDB(db)
	log.Println("Connected to database")
	return nil
}

func Disconnect() {
	sqlDB, _ := localDB.db.DB()
	defer sqlDB.Close()
	log.Println("Disconnected from database")
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
