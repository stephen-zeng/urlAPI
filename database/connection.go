package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"urlAPI/command"
)

func connect() {
	var err error
	os.Mkdir("assets", 0777)
	tmp, _ := sql.Open("sqlite3", dbPath)
	tmp.Close()
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Println(errors.Join(errors.New("database connection"), err))
		command.Exit()
	} else {
		log.Println("Connected to database")
	}
}

func Disconnect() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	log.Println("Disconnected from database")
}

func initMap() {
	var settings []Setting
	err := db.Where(1).Find(&settings).Error
	if err != nil {
		log.Println(errors.Join(errors.New("dataMap init"), err))
		command.Exit()
	}
	for _, setting := range settings {
		var settingList []string
		err = json.Unmarshal([]byte(setting.Value), &settingList)
		if err != nil {
			log.Println(errors.Join(errors.New("dataMap init"), err))
			command.Exit()
		}
		SettingMap[setting.Name] = settingList
	}
	log.Println("Initialized dataMap")
}

// 包括所有数据的初始化
func init() {
	connect()
	initMap()
}
