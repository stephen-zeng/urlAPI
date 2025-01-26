package data

import (
	"errors"
	"fmt"
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
		log.Fatal(err)
		fmt.Println("Error connecting to database")
		return errors.New("db.connect.error")
	}
	//defer db.Close()
	fmt.Println("Connected to database")
	return nil
}

func Disconnect() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	fmt.Println("Disconnected from database")
}

func init() {
	err := connect()
	if err == nil {
		taskInit()
		settingInit()
	}
}
