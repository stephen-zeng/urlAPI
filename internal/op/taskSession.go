package op

import (
	"github.com/pkg/errors"
	"reflect"
	"sort"
	"strings"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

var taskStatFields = map[string]string{
	"region":    "region",
	"type":      "type",
	"status":    "status",
	"api":       "api",
	"model":     "model",
	"referer":   "referer",
	"device":    "device",
	"more_info": "more_info",
	"temp":      "temp",
}

func fetchTask(info *Session) error {
	var taskGetter model.Task
	v := reflect.ValueOf(&taskGetter).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("json")
		if tag == info.TaskCatagory && tag != "time" {
			v.Field(i).Set(reflect.ValueOf(info.TaskBy))
		}
	}
	if info.TaskCatagory == "time" {
		taskGetter.Time = util.GetDate(info.TaskBy)
	}
	taskDBList, err := db.ReadTask(taskGetter)
	if err != nil {
		return errors.WithStack(err)
	}
	taskList := taskDBList.TaskList
	if len(taskList) == 0 {
		info.TaskMaxPage = 0
		info.TaskData = nil
		return nil
	}
	info.TaskMaxPage = ((len(taskList) - 1) / 100) + 1
	sort.Slice(taskList, func(i, j int) bool {
		return taskList[i].Time.After(taskList[j].Time)
	})
	switch {
	case info.TaskPage == -1:
		info.TaskData = taskList
	default:
		page := info.TaskPage
		if page < 1 {
			page = 1
		}
		start := (page - 1) * 100
		if start >= len(taskList) {
			info.TaskData = nil
			return nil
		}
		end := start + 100
		if end > len(taskList) {
			end = len(taskList)
		}
		info.TaskData = taskList[start:end]
	}
	return nil
}

func fetchTaskStats(info *Session) error {
	info.TaskStats = make(TaskStats, len(taskStatFields)+1)
	for key, field := range taskStatFields {
		stats, err := db.ReadTaskStats(field)
		if err != nil {
			return errors.WithStack(err)
		}
		info.TaskStats[key] = toTaskStatMap(stats)
	}

	timeStats, err := db.ReadTaskStats("strftime('%Y.%m', time)")
	if err != nil {
		return errors.WithStack(err)
	}
	info.TaskStats["time"] = toTaskStatMap(timeStats)
	return nil
}

func toTaskStatMap(stats []database.TaskStatItem) []TaskStatItem {
	ret := make([]TaskStatItem, 0, len(stats))
	for _, stat := range stats {
		ret = append(ret, TaskStatItem{
			Key:   strings.TrimSpace(stat.Key),
			Count: stat.Count,
		})
	}
	return ret
}
