package util

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseWebTarget(t *testing.T) {
	tests := []struct {
		raw     string
		wantErr bool
	}{
		{"https://github.com/a/b", false},
		{"http://www.bilibili.com/video/BV1xx411c7mD", false},
		{"  https://arxiv.org/abs/2101.00001  ", false},
		{"", true},
		{"   ", true},
		{"h", true},
		{"://", true},
		{"%zz", true},
		{"ftp://github.com/a/b", true},
		{"javascript:alert(1)", true},
		{"github.com/a/b", true},
		{"https://", true},
		{"https:///path", true},
	}
	for _, tt := range tests {
		_, err := ParseWebTarget(tt.raw)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseWebTarget(%q) err = %v, wantErr %v", tt.raw, err, tt.wantErr)
		}
	}
}

func TestWebTargetHost(t *testing.T) {
	tests := map[string]string{
		"https://GitHub.com/a/b":          "github.com",
		"https://www.youtube.com:443/x":   "www.youtube.com",
		"https://user:pw@gitee.com/a/b":   "gitee.com",
		"":                                "",
		"not a url":                       "",
		"https://[::1]:80/":               "::1",
		"https://www.bilibili.com/video/": "www.bilibili.com",
	}
	for raw, want := range tests {
		if got := WebTargetHost(raw); got != want {
			t.Errorf("WebTargetHost(%q) = %q, want %q", raw, got, want)
		}
	}
}

func TestBilibiliVideoID(t *testing.T) {
	tests := []struct {
		raw, want string
		wantErr   bool
	}{
		{"https://www.bilibili.com/video/BV1xx411c7mD", "BV1xx411c7mD", false},
		{"https://www.bilibili.com/video/BV1xx411c7mD/", "BV1xx411c7mD", false},
		{"https://www.bilibili.com/video/BV1xx411c7mD?p=2&t=10", "BV1xx411c7mD", false},
		{"https://www.bilibili.com/video/BV1xx411c7mD/?spm_id_from=333", "BV1xx411c7mD", false},
		{"https://www.bilibili.com/video/av170001", "av170001", false},
		{"https://www.bilibili.com/video/AV170001/", "AV170001", false},
		{"", "", true},
		{"https://www.bilibili.com", "", true},
		{"https://www.bilibili.com/video/", "", true},
		{"https://www.bilibili.com/video", "", true},
		{"https://www.bilibili.com/bangumi/play/ep1", "", true},
		{"https://www.bilibili.com/video/BV1", "", true},
		{"https://www.bilibili.com/video/av", "", true},
		{"https://www.bilibili.com/video/xyz123", "", true},
		{"https://www.bilibili.com/video/%2E%2E/BV1xx411c7mD", "", true},
		{"short", "", true},
	}
	for _, tt := range tests {
		got, err := BilibiliVideoID(tt.raw)
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("BilibiliVideoID(%q) = %q, %v; want %q, wantErr %v", tt.raw, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestYouTubeVideoID(t *testing.T) {
	tests := []struct {
		raw, want string
		wantErr   bool
	}{
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=42s", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/watch?list=PL1&v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/watch/?v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"", "", true},
		{"https://www.youtube.com/watch", "", true},
		{"https://www.youtube.com/watch?v=", "", true},
		{"https://www.youtube.com/watch?v=short", "", true},
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ%26key%3Dx", "", true},
		{"https://www.youtube.com/channel/abc", "", true},
		{"https://www.youtube.com", "", true},
		{"x", "", true},
	}
	for _, tt := range tests {
		got, err := YouTubeVideoID(tt.raw)
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("YouTubeVideoID(%q) = %q, %v; want %q, wantErr %v", tt.raw, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestArxivID(t *testing.T) {
	tests := []struct {
		raw, want string
		wantErr   bool
	}{
		{"https://arxiv.org/abs/2101.00001", "2101.00001", false},
		{"https://arxiv.org/abs/2101.00001v3", "2101.00001v3", false},
		{"https://arxiv.org/abs/1706.03762/", "1706.03762", false},
		{"https://arxiv.org/abs/1706.03762?context=cs", "1706.03762", false},
		{"https://arxiv.org/abs/hep-th/9901001", "hep-th/9901001", false},
		{"https://arxiv.org/abs/math.GT/0309136", "math.GT/0309136", false},
		{"", "", true},
		{"https://arxiv.org", "", true},
		{"https://arxiv.org/abs/", "", true},
		{"https://arxiv.org/pdf/2101.00001", "", true},
		{"https://arxiv.org/abs/../../etc/passwd", "", true},
		{"https://arxiv.org/abs/2101.00001/../x", "", true},
		{"https://arxiv.org/abs/hello", "", true},
		{"arxiv", "", true},
	}
	for _, tt := range tests {
		got, err := ArxivID(tt.raw)
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("ArxivID(%q) = %q, %v; want %q, wantErr %v", tt.raw, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestRepoFromURL(t *testing.T) {
	tests := []struct {
		raw, owner, repo string
		wantErr          bool
	}{
		{"https://github.com/stephen-zeng/urlAPI", "stephen-zeng", "urlAPI", false},
		{"https://github.com/stephen-zeng/urlAPI/", "stephen-zeng", "urlAPI", false},
		{"https://github.com/stephen-zeng/urlAPI.git", "stephen-zeng", "urlAPI", false},
		{"https://github.com/owner/repo/tree/main/src", "owner", "repo", false},
		{"https://gitee.com/owner/my.repo?x=1", "owner", "my.repo", false},
		{"", "", "", true},
		{"https://github.com", "", "", true},
		{"https://github.com/owner", "", "", true},
		{"https://github.com/owner/", "", "", true},
		{"https://github.com/../..", "", "", true},
		{"https://github.com/owner/%2e%2e", "", "", true},
		{"https://github.com/own%20er/repo", "", "", true},
		{"gh", "", "", true},
	}
	for _, tt := range tests {
		owner, repo, err := RepoFromURL(tt.raw)
		if (err != nil) != tt.wantErr || owner != tt.owner || repo != tt.repo {
			t.Errorf("RepoFromURL(%q) = %q, %q, %v; want %q, %q, wantErr %v", tt.raw, owner, repo, err, tt.owner, tt.repo, tt.wantErr)
		}
	}
}

func TestSplitRepoInfo(t *testing.T) {
	valid := []string{"owner/repo", "a-b/c_d.e"}
	invalid := []string{"", "/", "owner", "owner/", "/repo", "owner/repo/extra", "../x", "owner/..", "o w/r", "owner/repo?x"}
	for _, info := range valid {
		if _, _, err := SplitRepoInfo(info); err != nil {
			t.Errorf("SplitRepoInfo(%q) unexpected error %v", info, err)
		}
	}
	for _, info := range invalid {
		if _, _, err := SplitRepoInfo(info); err == nil {
			t.Errorf("SplitRepoInfo(%q) expected error", info)
		}
	}
}

func TestProvidersRejectInvalidIDsWithoutPanicking(t *testing.T) {
	if _, err := Bili(""); err == nil {
		t.Error("Bili(\"\") expected error")
	}
	if _, err := Bili("x"); err == nil {
		t.Error("Bili(\"x\") expected error")
	}
	if _, err := Ytb("", "token"); err == nil {
		t.Error("Ytb(\"\") expected error")
	}
	if _, err := Arxiv(""); err == nil {
		t.Error("Arxiv(\"\") expected error")
	}
	if _, err := Repo("", ""); err == nil {
		t.Error("Repo(\"\") expected error")
	}
	if _, err := Repo("https://example.com/a/b", ""); err == nil {
		t.Error("Repo on unsupported host expected error")
	}
	if _, err := ITHome("", "", "", "", ""); err == nil {
		t.Error("ITHome(\"\") expected error")
	}
	if _, err := ITHome("https://evil.example/0/1.htm", "", "", "", ""); err == nil {
		t.Error("ITHome on foreign host expected error")
	}
}

func TestYtbMalformedResponses(t *testing.T) {
	tests := map[string]string{
		"empty items":   `{"items":[]}`,
		"missing items": `{}`,
		"invalid json":  `{"items":`,
	}
	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(body))
			}))
			defer srv.Close()
			orig := youTubeVideoAPI
			youTubeVideoAPI = srv.URL
			defer func() { youTubeVideoAPI = orig }()

			if _, err := Ytb("dQw4w9WgXcQ", "secret-token"); err == nil {
				t.Fatal("expected error for malformed response")
			}
		})
	}
}

func TestRepoUsesCanonicalAPIPath(t *testing.T) {
	var gotPath, gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotAuth = r.URL.EscapedPath(), r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`not json`))
	}))
	defer srv.Close()
	orig := githubReposAPI
	githubReposAPI = srv.URL + "/repos/"
	defer func() { githubReposAPI = orig }()

	_, err := Repo("https://github.com/owner/repo.git/tree/main?tab=readme", "tok")
	if err == nil {
		t.Fatal("expected JSON decode error")
	}
	if gotPath != "/repos/owner/repo" {
		t.Fatalf("API path = %q, want /repos/owner/repo", gotPath)
	}
	if !strings.HasPrefix(gotAuth, "Bearer ") {
		t.Fatalf("GitHub token not forwarded: %q", gotAuth)
	}
}
