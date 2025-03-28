package database

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"reflect"
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
	var query string
	if !task.Time.IsZero() {
		start := task.Time
		end := start.AddDate(0, 1, 0)
		err = db.Where("time >= ? AND time <= ?", start, end).Find(&tasks).Error
	} else {
		val := reflect.ValueOf(task)
		if val.Kind() == reflect.Struct {
			for i := 0; i < val.NumField(); i++ {
				field := val.Type().Field(i).Tag.Get("json")
				value := val.Field(i)
				if value.IsNil() {
					continue
				}
				if value.Interface().(string) == "N/A" {
					query += fmt.Sprintf("(%s = %s OR %s IS NULL) AND ", field, "", field)
				} else {
					query += fmt.Sprintf("(%s = %s) AND ", field, value.Interface().(string))
				}
			}
		}
		if query == "" {
			err = db.Where(1).Find(&tasks).Error
		} else {
			err = db.Where(query[:len(query)-5]).Find(&tasks).Error
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
