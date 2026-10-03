package database

import (
	"fmt"
	"testing"
	"time"
	"urlAPI/internal/model"
)

func seedTasks(t *testing.T, adapter *SQLiteAdapter) {
	t.Helper()
	base := time.Date(2024, 5, 10, 12, 0, 0, 0, time.UTC)
	for i := 0; i < 250; i++ {
		task := model.Task{
			UUID:   fmt.Sprintf("task-%03d", i),
			Time:   base.Add(time.Duration(i) * time.Hour),
			Type:   []string{"txt", "img"}[i%2],
			Status: "success",
			API:    "openai",
			Region: "上海",
		}
		if i%10 == 0 {
			task.Region = ""
		}
		if err := adapter.CreateTask(&task); err != nil {
			t.Fatal(err)
		}
	}
	// A task in another month and one with a NULL region.
	other := model.Task{UUID: "april", Time: time.Date(2024, 4, 30, 23, 0, 0, 0, time.UTC), Type: "web", Region: "北京"}
	if err := adapter.CreateTask(&other); err != nil {
		t.Fatal(err)
	}
	if err := adapter.db.Exec("UPDATE tasks SET region = NULL WHERE uuid = ?", "task-005").Error; err != nil {
		t.Fatal(err)
	}
}

func TestQueryTasksPaginationAndOrdering(t *testing.T) {
	adapter := newTestDB(t)
	seedTasks(t, adapter)

	tasks, total, err := adapter.QueryTasks(TaskQuery{Page: 1, PageSize: 100})
	if err != nil {
		t.Fatal(err)
	}
	if total != 251 || len(tasks) != 100 {
		t.Fatalf("page 1: total=%d len=%d", total, len(tasks))
	}
	if tasks[0].UUID != "task-249" {
		t.Fatalf("newest first expected, got %s", tasks[0].UUID)
	}
	for i := 1; i < len(tasks); i++ {
		if tasks[i].Time.After(tasks[i-1].Time) {
			t.Fatalf("tasks not ordered by time DESC at %d", i)
		}
	}

	last, _, err := adapter.QueryTasks(TaskQuery{Page: 3, PageSize: 100})
	if err != nil || len(last) != 51 || last[len(last)-1].UUID != "april" {
		t.Fatalf("page 3: len=%d err=%v", len(last), err)
	}

	empty, total, err := adapter.QueryTasks(TaskQuery{Page: 4, PageSize: 100})
	if err != nil || len(empty) != 0 || total != 251 {
		t.Fatalf("page past the end: len=%d total=%d err=%v", len(empty), total, err)
	}

	all, _, err := adapter.QueryTasks(TaskQuery{Page: -1})
	if err != nil || len(all) != 251 {
		t.Fatalf("Page -1 should return everything: len=%d err=%v", len(all), err)
	}

	clamped, _, err := adapter.QueryTasks(TaskQuery{Page: 0, PageSize: 100})
	if err != nil || len(clamped) != 100 || clamped[0].UUID != "task-249" {
		t.Fatalf("Page 0 should behave as page 1: len=%d err=%v", len(clamped), err)
	}
}

func TestQueryTasksFilters(t *testing.T) {
	adapter := newTestDB(t)
	seedTasks(t, adapter)

	tests := []struct {
		name  string
		query TaskQuery
		want  int64
	}{
		{"type", TaskQuery{Field: "type", Value: "txt"}, 125},
		{"api", TaskQuery{Field: "api", Value: "openai"}, 250},
		{"no match", TaskQuery{Field: "status", Value: "failed"}, 0},
		{"N/A matches empty and NULL", TaskQuery{Field: "region", Value: TaskNotAvailable}, 26},
		{"region", TaskQuery{Field: "region", Value: "上海"}, 224},
		{"month", TaskQuery{Month: "2024.05"}, 250},
		{"other month", TaskQuery{Month: "2024.04"}, 1},
		{"empty month", TaskQuery{Month: "2023.01"}, 0},
		{"month and field", TaskQuery{Month: "2024.04", Field: "type", Value: "web"}, 1},
		{"value containing quotes", TaskQuery{Field: "target", Value: "x' OR '1'='1"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, total, err := adapter.QueryTasks(tt.query)
			if err != nil {
				t.Fatal(err)
			}
			if total != tt.want {
				t.Fatalf("total = %d, want %d", total, tt.want)
			}
			if int64(len(tasks)) > tt.want {
				t.Fatalf("returned %d tasks for %d matches", len(tasks), tt.want)
			}
		})
	}
}

func TestQueryTasksRejectsInvalidInput(t *testing.T) {
	adapter := newTestDB(t)
	invalid := []TaskQuery{
		{Month: "2024"},
		{Month: "2024.13"},
		{Month: "2024.00"},
		{Month: "24.05"},
		{Month: "2024.05; DROP TABLE tasks"},
		{Month: "."},
		{Field: "time", Value: "x"},
		{Field: "uuid = uuid OR 1", Value: "x"},
		{Field: "nonexistent", Value: "x"},
	}
	for _, q := range invalid {
		if _, _, err := adapter.QueryTasks(q); err == nil {
			t.Errorf("QueryTasks(%+v) expected error", q)
		}
	}
}

func TestMonthFilterMatchesStatsGrouping(t *testing.T) {
	adapter := newTestDB(t)
	shanghai := time.FixedZone("CST", 8*3600)
	times := []time.Time{
		time.Date(2024, 6, 1, 5, 0, 0, 0, shanghai), // 2024-05-31 21:00 UTC
		time.Date(2024, 6, 15, 12, 0, 0, 123456789, shanghai),
		time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
	}
	for i, ts := range times {
		if err := adapter.CreateTask(&model.Task{UUID: fmt.Sprint(i), Time: ts}); err != nil {
			t.Fatal(err)
		}
	}
	stats, err := adapter.ReadTaskStats("strftime('%Y.%m', time)")
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) == 0 {
		t.Fatal("no stats returned")
	}
	var sum int64
	for _, stat := range stats {
		_, total, err := adapter.QueryTasks(TaskQuery{Month: stat.Key})
		if err != nil {
			t.Fatalf("month %q: %v", stat.Key, err)
		}
		if total != int64(stat.Count) {
			t.Fatalf("month %q: filter total %d, stats count %d", stat.Key, total, stat.Count)
		}
		sum += total
	}
	if sum != int64(len(times)) {
		t.Fatalf("months cover %d tasks, want %d", sum, len(times))
	}
}
