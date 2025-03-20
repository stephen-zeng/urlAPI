package Database

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"time"
)

type Task struct {
	UUID    string    `json:"uuid" gorm:"primary_key"`
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

func taskInit() error {
	err := db.AutoMigrate(&Task{})
	if err != nil {
		return errors.Join(errors.New("Task Init"), err)
	}
	return nil
}

func (task *Task) Create() (*Task, error) {
	task.UUID = uuid.New().String()
	err := db.Create(task).Error
	if err != nil {
		return nil, errors.Join(errors.New("Task Create"), err)
	}
	return task, nil
}

func (task *Task) Update() (*Task, error) {
	err := db.Save(task).Error
	if err != nil {
		return nil, errors.Join(errors.New("Task Update"), err)
	}
	return task, nil
}

func (task *Task) Read() (*[]Task, error) {
	var ret []Task
	var err error
	var field string
	var value interface{}
	val := reflect.ValueOf(task)
	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			valField := val.Type().Field(i).Tag.Get("json")
			valValue := val.Field(i)
			if !valValue.IsNil() {
				field = valField
				value = valValue.Interface()
				break
			}
		}
	}
	switch field {
	case "":
		err = db.Where(1).Find(&ret).Error
	case "time":
		start := value.(time.Time)
		end := start.AddDate(0, 1, 0)
		err = db.Where("time >= ? AND time <= ?", start, end).Find(&ret).Error
	default:
		if value.(string) == "none" {
			err = db.Where("?=? OR ? IS NULL", field, "", field).Find(&ret).Error
		} else {
			err = db.Where("?=?", field, value.(string)).Find(&ret).Error
		}
	}
	if err != nil {
		return nil, errors.Join(errors.New("Task Read"), err)
	}
	return &ret, nil
}

func (task *Task) Delete() (*Task, error) {
	err := db.Delete(task).Error
	if err != nil {
		return nil, errors.Join(errors.New("Task Delete"), err)
	}
	return task, nil
}
