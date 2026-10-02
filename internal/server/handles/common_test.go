package handles

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlAPI/internal/op"
	"urlAPI/internal/server/trust"

	"github.com/gin-gonic/gin"
)

func runReturner(t *testing.T, cfg trust.Config, target, host string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	orig := trust.Get()
	trust.Set(cfg)
	t.Cleanup(func() { trust.Set(orig) })
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", target, nil)
	c.Request.Host = host
	c.Request.RemoteAddr = "203.0.113.10:5555"
	for k, v := range headers {
		c.Request.Header.Set(k, v)
	}
	returner(c, op.GenerateResult{URL: "/download?img=0f8fad5b-d9cb-469f-a165-70867728950e"})
	return w
}

func jsonURL(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var body op.GenerateResult
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v (%s)", err, w.Body.String())
	}
	return body.URL
}

func TestReturnerHostHandling(t *testing.T) {
	none, _ := trust.Parse("none", "", "")
	const path = "/download?img=0f8fad5b-d9cb-469f-a165-70867728950e"

	if got := jsonURL(t, runReturner(t, none, "/txt?format=json", "api.example.com", nil)); got != "http://api.example.com"+path {
		t.Fatalf("normal host: %q", got)
	}
	if got := jsonURL(t, runReturner(t, none, "/txt?format=json", "evil.example/x?", nil)); got != path {
		t.Fatalf("hostile host must not be embedded: %q", got)
	}
	// A spoofed X-Forwarded-Proto from an untrusted peer is ignored.
	spoofed := map[string]string{"X-Forwarded-Proto": "javascript"}
	if got := jsonURL(t, runReturner(t, none, "/txt?format=json", "api.example.com", spoofed)); got != "http://api.example.com"+path {
		t.Fatalf("spoofed proto: %q", got)
	}
	public, _ := trust.Parse("none", "https://api.example.com", "")
	if got := jsonURL(t, runReturner(t, public, "/txt?format=json", "evil.example", nil)); got != "https://api.example.com"+path {
		t.Fatalf("public origin: %q", got)
	}
	// Redirects use the relative URL, so the Host header is irrelevant.
	w := runReturner(t, none, "/txt", "evil.example", nil)
	if w.Code != http.StatusFound || w.Header().Get("Location") != path {
		t.Fatalf("redirect: %d %q", w.Code, w.Header().Get("Location"))
	}
}
