package trust

import (
	"crypto/tls"
	"net/http/httptest"
	"testing"
)

func mustParse(t *testing.T, proxies, public, dash string) Config {
	t.Helper()
	cfg, err := Parse(proxies, public, dash)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestParse(t *testing.T) {
	if cfg := mustParse(t, "", "", ""); len(cfg.Proxies) != len(DefaultTrustedProxies) {
		t.Fatalf("defaults not applied: %v", cfg.Proxies)
	}
	if cfg := mustParse(t, "none", "", ""); len(cfg.Proxies) != 0 {
		t.Fatalf("none should trust nothing: %v", cfg.Proxies)
	}
	cfg := mustParse(t, "203.0.113.7, 2001:db8::/32", "https://API.example.com/", "http://localhost:5173")
	if got := cfg.ProxyStrings(); len(got) != 2 || got[0] != "203.0.113.7/32" || got[1] != "2001:db8::/32" {
		t.Fatalf("ProxyStrings = %v", got)
	}
	if cfg.PublicOrigin != "https://api.example.com" {
		t.Fatalf("PublicOrigin = %q", cfg.PublicOrigin)
	}
	if !cfg.AllowsDashboardOrigin("http://localhost:5173") || cfg.AllowsDashboardOrigin("http://evil.example") {
		t.Fatal("dashboard origin matching is wrong")
	}
	bad := [][3]string{
		{"not-an-ip", "", ""},
		{"10.0.0.0/99", "", ""},
		{"", "ftp://example.com", ""},
		{"", "https://example.com/path", ""},
		{"", "https://user@example.com", ""},
		{"", "example.com", ""},
		{"", "", "*"},
	}
	for _, b := range bad {
		if _, err := Parse(b[0], b[1], b[2]); err == nil {
			t.Errorf("Parse(%q, %q, %q) expected error", b[0], b[1], b[2])
		}
	}
}

func TestSchemeHonoursForwardedProtoOnlyFromTrustedPeers(t *testing.T) {
	cfg := mustParse(t, "", "", "")
	tests := []struct {
		remote, proto, want string
	}{
		{"127.0.0.1:1234", "https", "https"},
		{"172.18.0.2:1234", "HTTPS", "https"},
		{"[::1]:1234", "https, http", "https"},
		{"127.0.0.1:1234", "http", "http"},
		{"127.0.0.1:1234", "javascript", "http"},
		{"127.0.0.1:1234", "https://evil.example/", "http"},
		{"203.0.113.9:1234", "https", "http"}, // spoofed by a public client
		{"[::ffff:127.0.0.1]:80", "https", "https"},
		{"garbage", "https", "http"},
	}
	for _, tt := range tests {
		r := httptest.NewRequest("GET", "/", nil)
		r.RemoteAddr = tt.remote
		r.Header.Set("X-Forwarded-Proto", tt.proto)
		if got := cfg.Scheme(r); got != tt.want {
			t.Errorf("Scheme(remote=%s, proto=%q) = %s, want %s", tt.remote, tt.proto, got, tt.want)
		}
	}
	r := httptest.NewRequest("GET", "/", nil)
	r.TLS = &tls.ConnectionState{}
	if cfg.Scheme(r) != "https" {
		t.Error("TLS request should be https")
	}
}

func TestOriginRejectsHostileHosts(t *testing.T) {
	cfg := mustParse(t, "none", "", "")
	tests := map[string]string{
		"api.example.com":      "http://api.example.com",
		"API.Example.com:8080": "http://api.example.com:8080",
		"127.0.0.1:2233":       "http://127.0.0.1:2233",
		"[::1]:2233":           "http://[::1]:2233",
		"evil.example/path":    "",
		"evil.example?x=1":     "",
		"user@evil.example":    "",
		"evil.example:abc":     "",
		"evil example":         "",
		"evil.example\\x":      "",
		"..":                   "",
		"":                     "",
		"[::1":                 "",
		"evil.example:":        "",
		"<script>":             "",
	}
	for host, want := range tests {
		r := httptest.NewRequest("GET", "/", nil)
		r.Host = host
		if got := cfg.Origin(r); got != want {
			t.Errorf("Origin(Host=%q) = %q, want %q", host, got, want)
		}
	}

	public := mustParse(t, "none", "https://api.example.com", "")
	r := httptest.NewRequest("GET", "/", nil)
	r.Host = "evil.example"
	if got := public.Absolute(r, "/download?img=x"); got != "https://api.example.com/download?img=x" {
		t.Fatalf("public origin ignored: %q", got)
	}
}

func TestAbsolute(t *testing.T) {
	cfg := mustParse(t, "none", "", "")
	r := httptest.NewRequest("GET", "/", nil)
	r.Host = "api.example.com"
	tests := map[string]string{
		"/download?img=x":           "http://api.example.com/download?img=x",
		"https://cdn.example.com/a": "https://cdn.example.com/a",
		"//evil.example/a":          "//evil.example/a",
		"":                          "",
		"download?img=legacy":       "download?img=legacy",
	}
	for in, want := range tests {
		if got := cfg.Absolute(r, in); got != want {
			t.Errorf("Absolute(%q) = %q, want %q", in, got, want)
		}
	}
	r.Host = "evil.example/x"
	if got := cfg.Absolute(r, "/download?img=x"); got != "/download?img=x" {
		t.Errorf("hostile Host should yield a relative URL, got %q", got)
	}
}
