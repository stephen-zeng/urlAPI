package database

import (
	"github.com/pkg/errors"
	"reflect"
	"urlAPI/internal/model"
)

type TaskStatItem struct {
	Key   string
	Count int
}

func (adapter *SQLiteAdapter) CreateTask(task *model.Task) error {
	return errors.WithStack(adapter.db.Create(task).Error)
}

func (adapter *SQLiteAdapter) UpdateTask(task *model.Task) error {
	return errors.WithStack(adapter.db.Save(task).Error)
}

func (adapter *SQLiteAdapter) ReadTask(task model.Task) (*model.DBList, error) {
	var tasks []model.Task
	var err error
	query := adapter.db.Model(&model.Task{})
	if !task.Time.IsZero() {
		start := task.Time
		end := start.AddDate(0, 1, 0)
		query.Where("time >= ? AND time <= ?", start, end)
	} else {
		val := reflect.ValueOf(task)
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
	ret := model.DBList{
		TaskList: tasks,
	}
	return &ret, errors.WithStack(err)
}

func (adapter *SQLiteAdapter) ReadTaskStats(field string) ([]TaskStatItem, error) {
	var stats []TaskStatItem
	expr := "COALESCE(" + field + ", '')"
	err := adapter.db.Model(&model.Task{}).
		Select(expr + " AS key, COUNT(*) AS count").
		Group("key").
		Order("count DESC").
		Find(&stats).Error
	return stats, errors.WithStack(err)
}

func (adapter *SQLiteAdapter) DeleteTask(task *model.Task) error {
	return errors.WithStack(adapter.db.Delete(task).Error)
}

func CreateTask(task *model.Task) error                  { return localDB.CreateTask(task) }
func UpdateTask(task *model.Task) error                  { return localDB.UpdateTask(task) }
func ReadTask(task model.Task) (*model.DBList, error)    { return localDB.ReadTask(task) }
func ReadTaskStats(field string) ([]TaskStatItem, error) { return localDB.ReadTaskStats(field) }
func DeleteTask(task *model.Task) error                  { return localDB.DeleteTask(task) }
