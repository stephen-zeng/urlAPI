package op

import (
	"fmt"
	"testing"
	"time"
	"urlAPI/internal/model"
)

func TestFetchTask(t *testing.T) {
	newTestDB(t)
	base := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 201; i++ {
		task := model.Task{UUID: fmt.Sprint(i), Time: base.Add(time.Duration(i) * time.Minute), Type: "txt"}
		if i == 0 {
			task.Type = ""
		}
		if err := db.CreateTask(&task); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name      string
		session   Session
		wantLen   int
		wantPages int
		wantErr   bool
	}{
		{"no filter", Session{TaskCatagory: "none", TaskPage: 1}, 100, 3, false},
		{"empty category", Session{TaskPage: 3}, 1, 3, false},
		{"past last page", Session{TaskCatagory: "none", TaskPage: 4}, 0, 3, false},
		{"all", Session{TaskCatagory: "none", TaskPage: -1}, 201, 3, false},
		{"type filter", Session{TaskCatagory: "type", TaskBy: "txt", TaskPage: 1}, 100, 2, false},
		{"N/A", Session{TaskCatagory: "type", TaskBy: "N/A", TaskPage: 1}, 1, 1, false},
		{"empty value means no filter", Session{TaskCatagory: "type", TaskPage: 1}, 100, 3, false},
		{"time", Session{TaskCatagory: "time", TaskBy: "2024.05", TaskPage: 1}, 100, 3, false},
		{"empty month", Session{TaskCatagory: "time", TaskBy: "2023.05", TaskPage: 1}, 0, 0, false},
		{"bad month", Session{TaskCatagory: "time", TaskBy: "2024", TaskPage: 1}, 0, 0, true},
		{"empty month value", Session{TaskCatagory: "time", TaskBy: "", TaskPage: 1}, 100, 3, false},
		{"unknown category", Session{TaskCatagory: "password", TaskBy: "x", TaskPage: 1}, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := tt.session
			err := fetchTask(&info)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if len(info.TaskData) != tt.wantLen || info.TaskMaxPage != tt.wantPages {
				t.Fatalf("len=%d pages=%d, want len=%d pages=%d", len(info.TaskData), info.TaskMaxPage, tt.wantLen, tt.wantPages)
			}
		})
	}
}
