package database

import (
	"database/sql"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"urlAPI/command"
)

var (
	dbPath string = "assets/database.db"
	db     *gorm.DB
)

func connect() error {
	var err error
	os.Mkdir("assets", 0777)
	tmp, _ := sql.Open("sqlite3", dbPath)
	tmp.Close()
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return errors.Join(errors.New("database connection"), err)
	} else {
		log.Println("Connected to database")
		return nil
	}
}

func Disconnect() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	log.Println("Disconnected from database")
}

// 包括所有数据的初始化
func init() {
	err := connect()
	if err != nil {
		log.Println(err)
		command.Exit()
	}

}
