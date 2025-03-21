package database

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"time"
)

func taskInit() error {
	err := db.AutoMigrate(&Task{})
	if err != nil {
		return errors.Join(errors.New("Task Init"), err)
	}
	return nil
}

func (task *Task) Create() error {
	task.UUID = uuid.New().String()
	err := db.Create(task).Error
	if err != nil {
		return errors.Join(errors.New("Task Create"), err)
	}
	return nil
}

func (task *Task) Update() error {
	err := db.Save(task).Error
	if err != nil {
		return errors.Join(errors.New("Task Update"), err)
	}
	return nil
}

func (task *Task) Read() (*DBList, error) {
	var tasks []Task
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
		err = db.Where(1).Find(&tasks).Error
	case "time":
		start := value.(time.Time)
		end := start.AddDate(0, 1, 0)
		err = db.Where("time >= ? AND time <= ?", start, end).Find(&tasks).Error
	default:
		if value.(string) == "none" {
			err = db.Where("?=? OR ? IS NULL", field, "", field).Find(&tasks).Error
		} else {
			err = db.Where("?=?", field, value.(string)).Find(&tasks).Error
		}
	}
	if err != nil {
		return nil, errors.Join(errors.New("Task Read"), err)
	}
	return &DBList{
		TaskList: tasks,
	}, nil
}

func (task *Task) Delete() error {
	err := db.Delete(task).Error
	if err != nil {
		return errors.Join(errors.New("Task Delete"), err)
	}
	return nil
}
