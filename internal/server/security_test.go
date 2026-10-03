package server

import (
	"encoding/json"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/internal/op"
	"urlAPI/internal/server/middleware"
	"urlAPI/internal/server/trust"
	"urlAPI/util"

	"github.com/gin-gonic/gin"
)

const testUUID = "0f8fad5b-d9cb-469f-a165-70867728950e"

func setupRouter(t *testing.T, cfg trust.Config, mutate func(*util.AppSettings)) *gin.Engine {
	t.Helper()
	origSettings, origImg, origTrust := database.SettingsStore.Get(), op.ImgPath, trust.Get()
	t.Cleanup(func() {
		database.SettingsStore.Replace(origSettings)
		op.ImgPath = origImg
		trust.Set(origTrust)
	})
	settings := util.AppSettings{}
	settings.Security.AllowedReferers = []string{"*"}
	settings.Security.DashboardAllowedIPs = []string{"*"}
	if mutate != nil {
		mutate(&settings)
	}
	database.SettingsStore.Replace(settings)

	op.ImgPath = t.TempDir() + "/"
	f, err := os.Create(op.ImgPath + testUUID + ".png")
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatal(err)
	}
	f.Close()

	r, err := NewRouter(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func defaultTrust(t *testing.T) trust.Config {
	t.Helper()
	cfg, err := trust.Parse("", "", "")
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

var remoteCounter int

// do sends a request from a fresh public IP so the rate limiter never
// interferes between cases.
func do(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	if req.RemoteAddr == "192.0.2.1:1234" { // httptest default
		remoteCounter++
		req.RemoteAddr = "198.51.100." + string(rune('0'+remoteCounter%10)) + ":1234"
		middleware.IPFrequency.IPFrequency = map[middleware.FrequencyFilter]middleware.FrequencyData{}
	}
	if req.Header.Get("Referer") == "" {
		req.Header.Set("Referer", "https://blog.example/post")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestDownloadAcceptsOnlyImageIDs(t *testing.T) {
	r := setupRouter(t, defaultTrust(t), nil)
	tests := []struct {
		query string
		want  int
	}{
		{"img=" + testUUID, http.StatusOK},
		{"img=" + strings.ToUpper(testUUID), http.StatusNotFound}, // valid form, file missing, no fallback
		{"img=empty", http.StatusOK},
		{"img=../../../../etc/passwd", http.StatusBadRequest},
		{"img=..%2F..%2Fetc%2Fpasswd", http.StatusBadRequest},
		{"img=%2e%2e%2f%2e%2e%2fdatabase", http.StatusBadRequest},
		{"img=" + testUUID + "/../" + testUUID, http.StatusBadRequest},
		{"img=" + testUUID + "%00", http.StatusBadRequest},
		{"img=..%5c..%5cwindows", http.StatusBadRequest},
		{"img=/etc/passwd", http.StatusBadRequest},
		{"img=EMPTY", http.StatusBadRequest},
		{"img=0f8fad5bd9cb469fa16570867728950e", http.StatusBadRequest},
		{"img={" + testUUID + "}", http.StatusBadRequest},
		{"img=urn:uuid:" + testUUID, http.StatusBadRequest},
		{"img=abc", http.StatusBadRequest},
	}
	for _, tt := range tests {
		w := do(r, httptest.NewRequest("GET", "/download?"+tt.query, nil))
		if tt.want == http.StatusBadRequest && w.Body.Len() > 0 && strings.Contains(w.Body.String(), "root:") {
			t.Fatalf("/download?%s leaked file contents", tt.query)
		}
		if w.Code != tt.want {
			t.Errorf("/download?%s: status %d, want %d (%s)", tt.query, w.Code, tt.want, w.Body.String())
		}
		if w.Code == http.StatusOK && w.Header().Get("Content-Type") != "image/png" {
			t.Errorf("/download?%s: content type %q", tt.query, w.Header().Get("Content-Type"))
		}
	}
}

func sessionRequest(body, token string) *http.Request {
	req := httptest.NewRequest("POST", "/session", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	return req
}

func TestDownloadMissingImageUsesFallback(t *testing.T) {
	r := setupRouter(t, defaultTrust(t), func(s *util.AppSettings) {
		s.Web.FallbackImageURL = "https://cdn.example.com/fallback.png"
	})
	missing := "11111111-2222-3333-4444-555555555555"
	w := do(r, httptest.NewRequest("GET", "/download?img="+missing, nil))
	if w.Code != http.StatusFound || w.Header().Get("Location") != "https://cdn.example.com/fallback.png" {
		t.Fatalf("missing image: %d %v", w.Code, w.Header())
	}
}

func TestUnknownSessionOperationFails(t *testing.T) {
	r := setupRouter(t, defaultTrust(t), nil)
	database.Sessions.Set(model.Session{Token: "valid-token", Expire: time.Now().Add(time.Hour)})
	t.Cleanup(func() { database.Sessions.Delete("valid-token") })

	w := do(r, sessionRequest(`{"operation":"dropDatabase"}`, "valid-token"))
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "unknown operation") {
		t.Fatalf("unknown op: %d %s", w.Code, w.Body.String())
	}
	w = do(r, sessionRequest(`{"operation":""}`, "valid-token"))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty op: %d %s", w.Code, w.Body.String())
	}
	w = do(r, sessionRequest(`{"operation":"dropDatabase"}`, "wrong-token"))
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "Authentication failed") {
		t.Fatalf("unauthenticated op: %d %s", w.Code, w.Body.String())
	}
}

func TestDashboardIPRestrictionIgnoresSpoofedHeaders(t *testing.T) {
	r := setupRouter(t, defaultTrust(t), func(s *util.AppSettings) {
		s.Security.DashboardAllowedIPs = []string{"10.1.2.3"}
	})
	tests := []struct {
		name, remote, xff string
		allowed           bool
	}{
		{"public client spoofing XFF", "203.0.113.5:4000", "10.1.2.3", false},
		{"public client spoofing X-Real-IP", "203.0.113.5:4000", "", false},
		{"trusted proxy forwarding allowed client", "127.0.0.1:4000", "10.1.2.3", true},
		{"trusted proxy forwarding other client", "127.0.0.1:4000", "203.0.113.5", false},
		{"allowed client directly", "10.1.2.3:4000", "", true},
	}
	for _, tt := range tests {
		req := sessionRequest(`{"operation":"noop"}`, "x")
		req.RemoteAddr = tt.remote
		if tt.xff != "" {
			req.Header.Set("X-Forwarded-For", tt.xff)
		} else if tt.name == "public client spoofing X-Real-IP" {
			req.Header.Set("X-Real-IP", "10.1.2.3")
		}
		w := do(r, req)
		denied := w.Code == http.StatusForbidden
		if denied == tt.allowed {
			t.Errorf("%s: status %d (%s)", tt.name, w.Code, w.Body.String())
		}
	}
}

func TestCORSPolicy(t *testing.T) {
	r := setupRouter(t, defaultTrust(t), nil)

	pre := httptest.NewRequest("OPTIONS", "/txt?prompt=x", nil)
	pre.Header.Set("Origin", "https://blog.example")
	pre.Header.Set("Access-Control-Request-Method", "GET")
	w := do(r, pre)
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("public API preflight: %d %v", w.Code, w.Header())
	}

	get := httptest.NewRequest("GET", "/download?img=empty", nil)
	get.Header.Set("Origin", "https://blog.example")
	w = do(r, get)
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("public API GET missing CORS header: %v", w.Header())
	}

	sess := sessionRequest(`{"operation":"login"}`, "x")
	sess.Header.Set("Origin", "https://evil.example")
	w = do(r, sess)
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "" {
		t.Fatalf("/session exposed to arbitrary origin: %q", v)
	}
	sessPre := httptest.NewRequest("OPTIONS", "/session", nil)
	sessPre.Header.Set("Origin", "https://evil.example")
	sessPre.Header.Set("Access-Control-Request-Method", "POST")
	w = do(r, sessPre)
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "" {
		t.Fatalf("/session preflight allowed arbitrary origin: %q", v)
	}

	cfg, err := trust.Parse("", "", "http://localhost:5173")
	if err != nil {
		t.Fatal(err)
	}
	r = setupRouter(t, cfg, nil)
	sessPre = httptest.NewRequest("OPTIONS", "/session", nil)
	sessPre.Header.Set("Origin", "http://localhost:5173")
	sessPre.Header.Set("Access-Control-Request-Method", "POST")
	sessPre.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
	w = do(r, sessPre)
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "http://localhost:5173" {
		t.Fatalf("configured dashboard origin rejected: %d %v", w.Code, w.Header())
	}
}

func TestRecoveryReturns500(t *testing.T) {
	r := setupRouter(t, defaultTrust(t), nil)
	r.GET("/panic-test", func(*gin.Context) { panic("boom") })
	w := do(r, httptest.NewRequest("GET", "/panic-test", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status %d", w.Code)
	}
	var v any
	_ = json.Unmarshal(w.Body.Bytes(), &v)
}
