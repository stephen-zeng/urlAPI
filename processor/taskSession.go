package processor

import (
	"github.com/pkg/errors"
	"reflect"
	"sort"
	"urlAPI/database"
)

func fetchTask(info *Session, data *database.Session) error {
	var taskGetter database.Task
	v := reflect.ValueOf(&taskGetter).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("json")
		if tag == info.TaskCatagory {
			v.Field(i).Set(reflect.ValueOf(info.TaskBy))
		}
	}
	taskDBList, err := taskGetter.Read()
	if err != nil {
		return errors.WithStack(err)
	}
	taskList := taskDBList.TaskList
	info.TaskMaxPage = (len(taskList) / 100) + 1
	sort.Slice(taskList, func(i, j int) bool {
		return taskList[i].Time.After(taskList[j].Time)
	})
	switch {
	case info.TaskPage == -1:
		info.TaskData = taskList
	case info.TaskPage*100 > len(taskList):
		info.TaskData = taskList[(info.TaskPage-1)*100:]
	case info.TaskPage*100 <= len(taskList):
		info.TaskData = taskList[(info.TaskPage-1)*100 : (info.TaskPage*100 - 1)]
	}
	return nil
}
