package database

import (
	"regexp"
	"strconv"
	"urlAPI/internal/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TaskStatItem struct {
	Key   string
	Count int
}

// TaskFilterColumns lists the task columns that may be used as equality
// filters. Keys are API field names; values are SQL column names. Only these
// identifiers are ever interpolated into SQL.
var TaskFilterColumns = map[string]string{
	"uuid":      "uuid",
	"ip":        "ip",
	"type":      "type",
	"status":    "status",
	"target":    "target",
	"return":    "return",
	"region":    "region",
	"referer":   "referer",
	"device":    "device",
	"more_info": "more_info",
	"api":       "api",
	"model":     "model",
	"temp":      "temp",
	"size":      "size",
}

var monthPattern = regexp.MustCompile(`^([0-9]{4})\.([0-9]{2})$`)

// TaskNotAvailable is the filter value matching empty or NULL columns.
const TaskNotAvailable = "N/A"

// TaskQuery describes a filtered, paginated task listing.
type TaskQuery struct {
	// Field is a key of TaskFilterColumns; empty means no column filter.
	Field string
	Value string
	// Month, formatted "YYYY.MM", restricts results to that calendar month
	// as reported by SQLite's strftime (the same grouping the stats use).
	Month string
	// Page is 1-based; a negative page returns every matching task.
	Page     int
	PageSize int
}

func (adapter *SQLiteAdapter) CreateTask(task *model.Task) error {
	return errors.WithStack(adapter.db.Create(task).Error)
}

func (adapter *SQLiteAdapter) UpdateTask(task *model.Task) error {
	return errors.WithStack(adapter.db.Save(task).Error)
}

func (adapter *SQLiteAdapter) filteredTasks(q TaskQuery) (*gorm.DB, error) {
	query := adapter.db.Model(&model.Task{})
	if q.Month != "" {
		m := monthPattern.FindStringSubmatch(q.Month)
		if m == nil {
			return nil, errors.Errorf("invalid month filter %q, expected YYYY.MM", q.Month)
		}
		if month, _ := strconv.Atoi(m[2]); month < 1 || month > 12 {
			return nil, errors.Errorf("invalid month filter %q, expected YYYY.MM", q.Month)
		}
		query = query.Where("strftime('%Y.%m', time) = ?", q.Month)
	}
	if q.Field != "" {
		column, ok := TaskFilterColumns[q.Field]
		if !ok {
			return nil, errors.Errorf("unsupported task filter %q", q.Field)
		}
		quoted := "`" + column + "`"
		if q.Value == TaskNotAvailable {
			query = query.Where("(" + quoted + " = '' OR " + quoted + " IS NULL)")
		} else {
			query = query.Where(quoted+" = ?", q.Value)
		}
	}
	return query, nil
}

// QueryTasks returns one page of matching tasks, newest first, together with
// the total number of matches.
func (adapter *SQLiteAdapter) QueryTasks(q TaskQuery) ([]model.Task, int64, error) {
	query, err := adapter.filteredTasks(q)
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	if total == 0 {
		return nil, 0, nil
	}
	list := query.Session(&gorm.Session{}).Order("time DESC")
	if q.Page >= 0 {
		page, size := q.Page, q.PageSize
		if page < 1 {
			page = 1
		}
		if size < 1 {
			size = 100
		}
		offset := int64(page-1) * int64(size)
		if offset >= total {
			return nil, total, nil
		}
		list = list.Limit(size).Offset(int(offset))
	}
	var tasks []model.Task
	if err := list.Find(&tasks).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return tasks, total, nil
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

func CreateTask(task *model.Task) error { return localDB.CreateTask(task) }
func UpdateTask(task *model.Task) error { return localDB.UpdateTask(task) }
func QueryTasks(q TaskQuery) ([]model.Task, int64, error) {
	return localDB.QueryTasks(q)
}
func ReadTaskStats(field string) ([]TaskStatItem, error) { return localDB.ReadTaskStats(field) }
func DeleteTask(task *model.Task) error                  { return localDB.DeleteTask(task) }
