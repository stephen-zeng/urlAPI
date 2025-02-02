package data

import (
	"database/sql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	dbPath string = "assets/database.db"
	db     *gorm.DB
	err    error
)

func connect() error {
	os.Mkdir("assets", 0777)
	tmp, _ := sql.Open("sqlite3", dbPath)
	tmp.Close()
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
		pwd, _ := InitSetting(DataConfig())
		InitSession(DataConfig())
		if pwd != "" {
			log.Printf("Dashboard password is %s\n", pwd)
		}
	}
}
