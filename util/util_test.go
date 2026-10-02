package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRegion(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch r.URL.Query().Get("ip") {
		case "1.1.1.1":
			_, _ = w.Write([]byte(`{"data":{"country":"中国","province":"上海"}}`))
		case "8.8.8.8":
			_, _ = w.Write([]byte(`{"data":{"country":"美国","province":""}}`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()
	origAPI, origCache := regionAPI, ipRegions
	regionAPI = srv.URL
	ipRegions = &regionCache{m: make(map[string]string)}
	t.Cleanup(func() { regionAPI, ipRegions = origAPI, origCache })

	if got := GetRegion("1.1.1.1"); got != "上海" {
		t.Fatalf("GetRegion(1.1.1.1) = %q", got)
	}
	if got := GetRegion("8.8.8.8"); got != "美国" {
		t.Fatalf("GetRegion(8.8.8.8) = %q", got)
	}
	if got := GetRegion("1.1.1.1"); got != "上海" || calls != 2 {
		t.Fatalf("cached lookup = %q after %d calls", got, calls)
	}
	if got := GetRegion("9.9.9.9"); got != "Unknown" {
		t.Fatalf("GetRegion on upstream failure = %q, want Unknown", got)
	}
}
