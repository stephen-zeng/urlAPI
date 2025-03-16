package data

import (
	"errors"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	UUID    string    `gorm:"primaryKey" json:"uuid"`
	Time    time.Time `json:"time"`
	IP      string    `json:"ip"`
	Type    string    `json:"type"`
	Status  string    `json:"status"`
	Target  string    `json:"target"`
	Return  string    `json:"return"`
	Size    string    `json:"size"`
	API     string    `json:"api"`
	Region  string    `json:"region"`
	Model   string    `json:"model"`
	Referer string    `json:"referer"`
	Device  string    `json:"device"`
}

func addTask(Time time.Time, API, Type, IP, Status, Target, Region, Size, Model, Referer, Device string) (string, error) {
	id := uuid.New().String()
	err := db.Create(Task{
		UUID:    id,
		API:     API,
		Time:    Time,
		IP:      IP,
		Type:    Type,
		Status:  Status,
		Target:  Target,
		Region:  Region,
		Size:    Size,
		Model:   Model,
		Referer: Referer,
		Device:  Device,
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
	if data == "-1" {
		return nil
	}
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
func parseDate(ori string) (int, int) {
	parts := strings.Split(ori, ".")
	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	return year, month
}

func fetchTask(by, data string) ([]Task, error) {
	var ret []Task
	var err *gorm.DB
	switch by {
	case "time":
		year, month := parseDate(data)
		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		err = db.Where("time >= ? AND time < ?", start, end).Find(&ret)
	case "none":
		err = db.Find(&ret)
	default:
		if data == "N/A" {
			err = db.Where(by+"=? OR "+by+" IS NULL", "").Find(&ret)
		} else {
			err = db.Where(by+"=?", data).Find(&ret)
		}
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
