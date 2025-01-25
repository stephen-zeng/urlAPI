package data

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

type Task struct {
	UUID   string `gorm:"primaryKey"`
	Time   time.Time
	IP     string
	Type   string
	Status string
	Target string
	Return string
}

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

func dbInit() error {
	var err error
	if !db.Migrator().HasTable(&Task{}) {
		err = db.AutoMigrate(&Task{})
	}
	if err != nil {
		fmt.Println("Error creating table")
		log.Fatal(err)
		return errors.New("db.init.error")
	}
	fmt.Println("Initialized database")
	return nil
}

func Disconnect() {
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	fmt.Println("Disconnected from database")
}

func add(Time time.Time, IP, Type, Status, Target string) (string, error) {
	id := uuid.New().String()
	err := db.Create(Task{
		UUID:   id,
		Time:   Time,
		IP:     IP,
		Type:   Type,
		Status: Status,
		Target: Target,
	})
	if err.Error != nil {
		return "", errors.New("db.add.error")
	} else {
		return id, nil
	}
}

// by中需要是SQL里面的数据类型
// data中是by对应的值
func del(by, data string) error {
	err := db.Where(by+"=?", data).Delete(&Task{})
	if err.Error != nil {
		return errors.New("db.del.error")
	} else {
		return nil
	}
}

// by中需要是SQL里面的数据类型
// data中是by对应的值
// 现在只有更改Status和Return的需要
func edit(by, data string, Status, Return string) error {
	err := db.Model(&Task{}).Where(by+"=?", data).Updates(Task{
		Status: Status,
		Return: Return,
	})
	if err.Error != nil {
		return errors.New("db.edit.error")
	} else {
		return nil
	}
}

// by中需要是SQL里面的数据类型
// data中是by对应的值
func fetch(by, data string) ([]Task, error) {
	var ret []Task
	var err *gorm.DB
	if by == "none" {
		err = db.Find(&ret)
	} else {
		err = db.Where(by+"=?", data).Find(&ret)
	}
	if err.Error != nil {
		return nil, errors.New("db.fetch.error")
	}
	if len(ret) == 0 {
		return nil, errors.New("db.fetch.nodata")
	} else {
		return ret, nil
	}
}
