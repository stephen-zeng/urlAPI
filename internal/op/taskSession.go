package op

import (
	"github.com/pkg/errors"
	"strings"
	"urlAPI/internal/database"
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

// taskPageSize is the number of tasks per dashboard page.
const taskPageSize = 100

func fetchTask(info *Session) error {
	query := database.TaskQuery{Page: info.TaskPage, PageSize: taskPageSize}
	if info.TaskPage == -1 {
		query.Page = -1
	} else if query.Page < 1 {
		query.Page = 1
	}
	switch info.TaskCatagory {
	case "", "none":
	case "time":
		query.Month = info.TaskBy
	default:
		if _, ok := database.TaskFilterColumns[info.TaskCatagory]; !ok {
			return errors.Errorf("unsupported task filter %q", info.TaskCatagory)
		}
		// An empty value means "no filter", as before.
		if info.TaskBy != "" {
			query.Field = info.TaskCatagory
			query.Value = info.TaskBy
		}
	}
	tasks, total, err := db.QueryTasks(query)
	if err != nil {
		return errors.WithStack(err)
	}
	info.TaskMaxPage = int((total + taskPageSize - 1) / taskPageSize)
	info.TaskData = tasks
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
