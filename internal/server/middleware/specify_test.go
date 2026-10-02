package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"urlAPI/internal/database"
	"urlAPI/util"

	"github.com/gin-gonic/gin"
)

func webTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	orig := database.SettingsStore.Get()
	t.Cleanup(func() { database.SettingsStore.Replace(orig) })

	settings := util.AppSettings{}
	settings.Security.AllowedReferers = []string{"*"}
	settings.Features.WebImgEnabled = true
	settings.Web.AllowedHosts = []string{"github.com", "www.bilibili.com"}
	database.SettingsStore.Replace(settings)

	r := gin.New()
	r.GET("/web", GeneralSecurityMiddleware("web"), WebSecurityMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}

func TestWebSecurityMiddlewareRejectsMalformedTargets(t *testing.T) {
	r := webTestRouter(t)
	tests := []struct {
		target string
		want   int
	}{
		{"https://github.com/owner/repo", http.StatusOK},
		{"https://GitHub.com/owner/repo", http.StatusOK},
		{"%zz", http.StatusForbidden},
		{"://", http.StatusForbidden},
		{"h", http.StatusForbidden},
		{"ftp://github.com/owner/repo", http.StatusForbidden},
		{"https://evil.example/owner/repo", http.StatusForbidden},
	}
	for i, tt := range tests {
		resetFrequency(t)
		req := httptest.NewRequest(http.MethodGet, "/web?img="+url.QueryEscape(tt.target), nil)
		req.Header.Set("Referer", "https://site.example/")
		req.RemoteAddr = "10.0.0." + string(rune('1'+i)) + ":1234"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tt.want {
			t.Errorf("target %q: status %d, want %d (%s)", tt.target, w.Code, tt.want, w.Body.String())
		}
	}
}
