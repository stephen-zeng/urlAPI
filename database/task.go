package database

import (
	"errors"
	"reflect"
)

func (task *Task) Create() error {
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
	query := db.Model(&Task{})
	if !task.Time.IsZero() {
		start := task.Time
		end := start.AddDate(0, 1, 0)
		query.Where("time >= ? AND time <= ?", start, end)
	} else {
		val := reflect.ValueOf(*task)
		if val.Kind() == reflect.Struct {
			for i := 0; i < val.NumField(); i++ {
				field := val.Type().Field(i).Tag.Get("json")
				value := val.Field(i)
				if value.IsZero() {
					continue
				}
				if value.Interface().(string) == "N/A" {
					query.Where(field+"=? OR "+field+" IS NULL", "")
				} else {
					query.Where(field+"=?", value.Interface().(string))
				}
			}
		}
	}
	err = query.Find(&tasks).Error
	ret := DBList{
		TaskList: tasks,
	}
	if err != nil {
		return &ret, errors.Join(errors.New("Task Read"), err)
	}
	return &ret, nil
}

func (task *Task) Delete() error {
	err := db.Delete(task).Error
	if err != nil {
		return errors.Join(errors.New("Task Delete"), err)
	}
	return nil
}
