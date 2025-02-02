package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var (
	dbPath string = "assets/database.db"
	db     *gorm.DB
	err    error
)

func connect() error {
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return err
	}
	//defer db.Close()
	log.Println("Connected to database")
	return nil
}

func Disconnect() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	log.Println("Disconnected from database")
}

func init() {
	err := connect()
	if err == nil {
		initTask()
		InitSetting(DataConfig())
		InitSession(DataConfig())
	}
}
