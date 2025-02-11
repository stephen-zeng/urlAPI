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
	UUID   string    `gorm:"primaryKey" json:"uuid"`
	Time   time.Time `json:"time"`
	IP     string    `json:"ip"`
	Type   string    `json:"type"`
	Status string    `json:"status"`
	Target string    `json:"target"`
	Return string    `json:"return"`
	Size   string    `json:"size"`
	API    string    `json:"api"`
	Region string    `json:"region"`
}

func addTask(Time time.Time, API, Type, IP, Status, Target, Region, Size string) (string, error) {
	id := uuid.New().String()
	err := db.Create(Task{
		UUID:   id,
		API:    API,
		Time:   Time,
		IP:     IP,
		Type:   Type,
		Status: Status,
		Target: Target,
		Region: Region,
		Size:   Size,
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
func editTask(by, data string, Status, Return, Size, API, Region string) error {
	err := db.Model(&Task{}).Where(by+"=?", data).Updates(Task{
		Status: Status,
		Return: Return,
		Size:   Size,
		API:    API,
		Region: Region,
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
		return nil, errors.New("Task not found")
	} else {
		return ret, nil
	}
}
