package data

import (
	"errors"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
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
	Size   string
	API    string
}

func initTask() error {
	var err error
	if !db.Migrator().HasTable(&Task{}) {
		err = db.AutoMigrate(&Task{})
	}
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Initialized tasks table")
	return nil
}

func addTask(Time time.Time, IP, Type, Status, Target string) (string, error) {
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
		log.Println(err.Error)
		return "", err.Error
	} else {
		return id, nil
	}
}

// by中需要是SQL里面的数据类型
// data中是by对应的值
func delTask(by, data string) error {
	err := db.Where(by+"=?", data).Delete(&Task{})
	if err.Error != nil {
		log.Println(err.Error)
		return err.Error
	} else {
		return nil
	}
}

// by中需要是SQL里面的数据类型
// data中是by对应的值
// 现在只有更改Status和Return的需要
func editTask(by, data string, Status, Return, Size, API string) error {
	err := db.Model(&Task{}).Where(by+"=?", data).Updates(Task{
		Status: Status,
		Return: Return,
		Size:   Size,
		API:    API,
	})
	if err.Error != nil {
		return err.Error
	} else {
		return nil
	}
}

// by中需要是SQL里面的数据类型
// data中是by对应的值
func fetchTask(by, data string) ([]Task, error) {
	var ret []Task
	var err *gorm.DB
	if by == "none" {
		err = db.Find(&ret)
	} else {
		err = db.Where(by+"=?", data).Find(&ret)
	}
	if err.Error != nil {
		return nil, err.Error
	}
	if len(ret) == 0 {
		return nil, errors.New("record not found")
	} else {
		return ret, nil
	}
}
