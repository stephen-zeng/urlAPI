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
func init() {
	figlet := figure.NewFigure("urlAPI", "", true)
	figlet.Print()
	connect()
	migration()
	initSettingMap()
	initRepoMap()
	initSessionMap()
	if err := settingInit(); err != nil {
		log.Fatal(errors.Wrap(err, "settingInit"))
	}
}

func migration() {
	db.AutoMigrate(&Setting{})
	db.AutoMigrate(&Task{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Repo{})
}

func connect() {
	var err error
	os.Mkdir("assets", 0777)
	tmp, _ := sql.Open("sqlite3", dbPath)
	tmp.Close()
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal(errors.Wrap(err, "gorm"))
	} else {
		log.Println("Connected to database")
	}
}

func Disconnect() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	log.Println("Disconnected from database")
}

func initSettingMap() {
	var settings []Setting
	if err := db.Find(&settings).Error; err != nil {
		log.Fatal(errors.Wrap(err, "db"))
	}
	for _, setting := range settings {
		var settingList []string
		if err := json.Unmarshal([]byte(setting.Value), &settingList); err != nil {
			log.Fatal(errors.Wrap(err, "json"))
		}
		SettingMap[setting.Name] = settingList
	}
	log.Println("Initialized SettingMap")
}

func initRepoMap() {
	var repos []Repo
	if err := db.Find(&repos).Error; err != nil {
		log.Fatal(errors.Wrap(err, "db find").Error())
	}
	for _, repo := range repos {
		var repoList []string
		if err := json.Unmarshal([]byte(repo.Content), &repoList); err != nil {
			log.Fatal(errors.Wrap(err, "json"))
		}
		RepoMap[repo.API+";"+repo.Info] = repoList
	}
	log.Println("Initialized RepoMap")
}

func initSessionMap() {
	var sessions []Session
	if err := db.Find(&sessions).Error; err != nil {
		log.Fatal(errors.Wrap(err, "db"))
	}
	for _, session := range sessions {
		SessionMap[session.Token] = session
	}
	log.Println("Initialized SessionMap")
}

func ClearTask() {
	if db.Migrator().HasTable(&Task{}) {
		if err := db.Migrator().DropTable(&Task{}); err != nil {
			log.Fatal(errors.Wrap(err, "db"))
		}
		if err := db.AutoMigrate(&Task{}); err != nil {
			log.Fatal(errors.Wrap(err, "db"))
		}
	}
}

func ClearSession() {
	if db.Migrator().HasTable(&Session{}) {
		if err := db.Migrator().DropTable(&Session{}); err != nil {
			log.Fatal(errors.Wrap(err, "db"))
		}
		if err := db.AutoMigrate(&Session{}); err != nil {
			log.Fatal(errors.Wrap(err, "db"))
		}
	}
	SessionMap = make(map[string]Session)
}
