package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func connect() {
	var err error
	os.Mkdir("assets", 0777)
	tmp, _ := sql.Open("sqlite3", dbPath)
	tmp.Close()
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal(errors.Join(errors.New("database connection"), err))
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
	err := db.Where(1).Find(&settings).Error
	if err != nil {
		log.Fatal(errors.Join(errors.New("SettingMap init"), err))
	}
	for _, setting := range settings {
		var settingList []string
		err = json.Unmarshal([]byte(setting.Value), &settingList)
		if err != nil {
			log.Fatal(errors.Join(errors.New("SettingMap init"), err))
		}
		SettingMap[setting.Name] = settingList
	}
	log.Println("Initialized SettingMap")
}

func initRepoMap() {
	var repos []Repo
	err := db.Where(1).Find(&repos).Error
	if err != nil {
		log.Fatal(errors.Join(errors.New("RepoMap init"), err))
	}
	for _, repo := range repos {
		var repoList []string
		err = json.Unmarshal([]byte(repo.Content), &repoList)
		if err != nil {
			log.Fatal(errors.Join(errors.New("RepoMap init"), err))
		}
		RepoMap[repo.API+";"+repo.Info] = repoList
	}
	log.Println("Initialized RepoMap")
}

func initSessionMap() {
	var sessions []Session
	err := db.Where(1).Find(&sessions).Error
	if err != nil {
		log.Fatal(errors.Join(errors.New("SessionMap init"), err))
	}
	for _, session := range sessions {
		SessionMap[session.Token] = session
	}
	log.Println("Initialized SessionMap")
}

// 包括所有数据的初始化
func init() {
	connect()
	initSettingMap()
	initRepoMap()
	initSessionMap()
	if err := settingInit(); err != nil {
		log.Fatal(errors.Join(errors.New("SettingInit"), err))
	}
	if err := taskInit(); err != nil {
		log.Fatal(errors.Join(errors.New("TaskInit"), err))
	}
	if err := sessionInit(); err != nil {
		log.Fatal(errors.Join(errors.New("SessionInit"), err))
	}
	if err := repoInit(); err != nil {
		log.Fatal(errors.Join(errors.New("RepoInit"), err))
	}
}

func ClearTask() {
	if db.Migrator().HasTable(&Task{}) {
		if err := db.Migrator().DropTable(&Task{}); err != nil {
			log.Fatal(errors.Join(errors.New("ClearTask"), err))
		}
		if err := db.AutoMigrate(&Task{}); err != nil {
			log.Fatal(errors.Join(errors.New("ClearTask"), err))
		}
	}
}

func ClearSession() {
	if db.Migrator().HasTable(&Session{}) {
		if err := db.Migrator().DropTable(&Session{}); err != nil {
			log.Fatal(errors.Join(errors.New("ClearSession"), err))
		}
		if err := db.AutoMigrate(&Session{}); err != nil {
			log.Fatal(errors.Join(errors.New("ClearSession"), err))
		}
	}
}
func ClearSetting() {
	if db.Migrator().HasTable(&Setting{}) {
		if err := db.Migrator().DropTable(&Setting{}); err != nil {
			log.Fatal(errors.Join(errors.New("ClearSetting"), err))
		}
		if err := db.AutoMigrate(&Setting{}); err != nil {
			log.Fatal(errors.Join(errors.New("ClearSetting"), err))
		}
	}
}
