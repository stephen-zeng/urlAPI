package database

import (
	"github.com/pkg/errors"
	"reflect"
)

func (task *Task) Create() error {
	return errors.WithStack(db.Create(task).Error)
}

func (task *Task) Update() error {
	return errors.WithStack(db.Save(task).Error)
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
	return &ret, errors.WithStack(err)
}

func (task *Task) Delete() error {
	return errors.WithStack(db.Delete(task).Error)
}
