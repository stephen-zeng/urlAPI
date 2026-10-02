package util

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const testSecret = "sk-test-secret-123456"

func newServer(t *testing.T, h http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv
}

func assertNoSecret(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error")
	}
	if strings.Contains(err.Error(), testSecret) {
		t.Fatalf("error leaks secret: %v", err)
	}
}

func TestDoUpstreamSuccess(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+testSecret {
			t.Errorf("missing auth header")
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	var out struct{ OK bool }
	err := doUpstreamJSON(context.Background(), upstreamRequest{URL: srv.URL, Headers: bearer(testSecret)}, &out)
	if err != nil || !out.OK {
		t.Fatalf("doUpstreamJSON = %v, %+v", err, out)
	}
}

func TestDoUpstreamNon2xx(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"Incorrect API key provided: ` + testSecret + `"}}`))
	})
	_, err := doUpstream(context.Background(), upstreamRequest{URL: srv.URL, Secrets: []string{testSecret}})
	assertNoSecret(t, err)
	var upstream *UpstreamError
	if !errors.As(err, &upstream) || upstream.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected UpstreamError 401, got %v", err)
	}
	if !strings.Contains(err.Error(), "Incorrect API key") || !strings.Contains(err.Error(), "[REDACTED]") {
		t.Fatalf("error lacks redacted upstream detail: %v", err)
	}
}

func TestDoUpstreamNon2xxWithoutBody(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})
	_, err := doUpstream(context.Background(), upstreamRequest{URL: srv.URL})
	if err == nil || !strings.Contains(err.Error(), "502") {
		t.Fatalf("expected 502 error, got %v", err)
	}
}

func TestDoUpstreamMalformedJSON(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"choices": [`))
	})
	var out TxtResp
	if err := doUpstreamJSON(context.Background(), upstreamRequest{URL: srv.URL}, &out); err == nil {
		t.Fatal("expected decode error")
	}
}

func TestDoUpstreamTransportErrorRedactsURL(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	addr := srv.URL
	srv.Close() // connection refused from now on
	_, err := doUpstream(context.Background(), upstreamRequest{URL: addr + "/v3/videos?id=x&key=" + testSecret, Secrets: []string{testSecret}})
	assertNoSecret(t, err)
	if strings.Contains(err.Error(), "key=") {
		t.Fatalf("error leaks query string: %v", err)
	}
}

func TestDoUpstreamInvalidURL(t *testing.T) {
	_, err := doUpstream(context.Background(), upstreamRequest{URL: "http://[::1]:namedport/?key=" + testSecret})
	assertNoSecret(t, err)
}

func TestDoUpstreamBodyLimit(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("a", 2048)))
	})
	_, err := doUpstream(context.Background(), upstreamRequest{URL: srv.URL, Limit: 1024})
	if !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("expected ErrResponseTooLarge, got %v", err)
	}
	data, err := doUpstream(context.Background(), upstreamRequest{URL: srv.URL, Limit: 2048})
	if err != nil || len(data) != 2048 {
		t.Fatalf("exact-limit body: len=%d err=%v", len(data), err)
	}
}

func TestDoUpstreamTimeoutAndCancellation(t *testing.T) {
	release := make(chan struct{})
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	})
	defer close(release)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := doUpstream(ctx, upstreamRequest{URL: srv.URL})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if time.Since(start) > 5*time.Second {
		t.Fatal("request was not bounded by the context deadline")
	}

	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() { time.Sleep(20 * time.Millisecond); cancel2() }()
	if _, err := doUpstream(ctx2, upstreamRequest{URL: srv.URL}); err == nil {
		t.Fatal("expected cancellation error")
	}
}

func TestTxt(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		want    string
		wantErr bool
	}{
		{"success", 200, `{"choices":[{"message":{"role":"assistant","content":"hi"}}]}`, "hi", false},
		{"empty choices", 200, `{"choices":[]}`, "", true},
		{"malformed", 200, `nope`, "", true},
		{"server error", 500, `{"error":"boom"}`, "", true},
		{"created", 201, `{"choices":[{"message":{"content":"ok"}}]}`, "ok", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
				var payload TxtPayload
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || len(payload.Messages) != 2 {
					t.Errorf("bad payload: %v %+v", err, payload)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			})
			got, err := Txt(srv.URL, testSecret, "model", "system", "prompt")
			if (err != nil) != tt.wantErr || got != tt.want {
				t.Fatalf("Txt = %q, %v; want %q, wantErr %v", got, err, tt.want, tt.wantErr)
			}
			if err != nil {
				assertNoSecret(t, err)
			}
		})
	}
	if _, err := Txt("", "", "", "", ""); err == nil {
		t.Fatal("expected error for missing parameters")
	}
}

func TestOpenaiImg(t *testing.T) {
	var imgURL string
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			_, _ = w.Write([]byte(`{"data":[{"url":"` + imgURL + `"}]}`))
		case "/empty":
			_, _ = w.Write([]byte(`{"data":[]}`))
		case "/img":
			_, _ = w.Write([]byte("PNGDATA"))
		case "/missing-img":
			_, _ = w.Write([]byte(`{"data":[{"url":"` + strings.TrimSuffix(imgURL, "/img") + `/404"}]}`))
		default:
			http.NotFound(w, r)
		}
	})
	imgURL = srv.URL + "/img"

	data, err := OpenaiImg(srv.URL+"/ok", testSecret, "p", "m", "1x1")
	if err != nil || string(data) != "PNGDATA" {
		t.Fatalf("OpenaiImg = %q, %v", data, err)
	}
	if _, err := OpenaiImg(srv.URL+"/empty", testSecret, "p", "m", "1x1"); err == nil {
		t.Fatal("expected error for empty data")
	}
	if _, err := OpenaiImg(srv.URL+"/missing-img", testSecret, "p", "m", "1x1"); err == nil {
		t.Fatal("expected error when the image download fails")
	}
	if _, err := OpenaiImg(srv.URL+"/nope", testSecret, "p", "m", "1x1"); err == nil {
		t.Fatal("expected error for 404 endpoint")
	}
}

func setAlibabaTestServer(t *testing.T, srv *httptest.Server) {
	t.Helper()
	origSubmit, origTask, origInterval := alibabaImageSynthesisAPI, alibabaTaskAPI, alibabaPollInterval
	alibabaImageSynthesisAPI = srv.URL + "/submit"
	alibabaTaskAPI = srv.URL + "/tasks/"
	alibabaPollInterval = 5 * time.Millisecond
	t.Cleanup(func() {
		alibabaImageSynthesisAPI, alibabaTaskAPI, alibabaPollInterval = origSubmit, origTask, origInterval
	})
}

func TestAlibabaImgSuccess(t *testing.T) {
	var polls atomic.Int32
	var srvURL string
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/submit":
			if r.Header.Get("X-DashScope-Async") != "enable" {
				t.Errorf("missing async header")
			}
			_, _ = w.Write([]byte(`{"output":{"task_id":"t1","task_status":"PENDING"}}`))
		case r.URL.Path == "/tasks/t1":
			if polls.Add(1) < 3 {
				_, _ = w.Write([]byte(`{"output":{"task_id":"t1","task_status":"RUNNING"}}`))
				return
			}
			_, _ = w.Write([]byte(`{"output":{"task_id":"t1","task_status":"SUCCEEDED","results":[{"actual_prompt":"actual","url":"` + srvURL + `/img"}]}}`))
		case r.URL.Path == "/img":
			_, _ = w.Write([]byte("IMG"))
		default:
			http.NotFound(w, r)
		}
	})
	srvURL = srv.URL
	setAlibabaTestServer(t, srv)

	data, prompt, err := AlibabaImg(testSecret, "p", "m", "1*1")
	if err != nil || string(data) != "IMG" || prompt != "actual" {
		t.Fatalf("AlibabaImg = %q, %q, %v", data, prompt, err)
	}
}

func TestAlibabaImgFailures(t *testing.T) {
	tests := map[string]string{
		"failed status": `{"output":{"task_id":"t1","task_status":"FAILED"}}`,
		"no results":    `{"output":{"task_id":"t1","task_status":"SUCCEEDED","results":[]}}`,
		"malformed":     `{"output":`,
	}
	for name, pollBody := range tests {
		t.Run(name, func(t *testing.T) {
			srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/submit" {
					_, _ = w.Write([]byte(`{"output":{"task_id":"t1","task_status":"PENDING"}}`))
					return
				}
				_, _ = w.Write([]byte(pollBody))
			})
			setAlibabaTestServer(t, srv)
			_, _, err := AlibabaImg(testSecret, "p", "m", "1*1")
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoSecret(t, err)
		})
	}

	t.Run("submit rejected", func(t *testing.T) {
		srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"code":"InvalidParameter","message":"bad size"}`))
		})
		setAlibabaTestServer(t, srv)
		_, _, err := AlibabaImg(testSecret, "p", "m", "bad")
		if err == nil || !strings.Contains(err.Error(), "bad size") {
			t.Fatalf("expected upstream detail, got %v", err)
		}
	})

	t.Run("missing task id", func(t *testing.T) {
		srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"output":{}}`))
		})
		setAlibabaTestServer(t, srv)
		if _, _, err := AlibabaImg(testSecret, "p", "m", "1*1"); err == nil {
			t.Fatal("expected error for missing task id")
		}
	})
}

func TestAlibabaImgPollingTimeoutTerminates(t *testing.T) {
	var polls atomic.Int32
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/submit" {
			_, _ = w.Write([]byte(`{"output":{"task_id":"t1","task_status":"PENDING"}}`))
			return
		}
		polls.Add(1)
		_, _ = w.Write([]byte(`{"output":{"task_id":"t1","task_status":"RUNNING"}}`))
	})
	setAlibabaTestServer(t, srv)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, _, err := alibabaImg(ctx, testSecret, "p", "m", "1*1")
	if !errors.Is(err, context.DeadlineExceeded) && (err == nil || !strings.Contains(err.Error(), "deadline")) {
		t.Fatalf("expected deadline error, got %v", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("polling did not stop promptly: %v", elapsed)
	}
	after := polls.Load()
	time.Sleep(50 * time.Millisecond)
	if polls.Load() != after {
		t.Fatal("polling continued after the call returned")
	}
}

func TestDownloader(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ok" {
			_, _ = w.Write([]byte("data"))
			return
		}
		w.WriteHeader(http.StatusForbidden)
	})
	if data, err := Downloader(srv.URL + "/ok"); err != nil || string(data) != "data" {
		t.Fatalf("Downloader = %q, %v", data, err)
	}
	if _, err := Downloader(srv.URL + "/denied"); err == nil {
		t.Fatal("expected error for 403")
	}
	if _, err := Downloader("http://127.0.0.1:1/unreachable"); err == nil {
		t.Fatal("expected error for unreachable host")
	}
}

func TestGetRepo(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			_, _ = w.Write([]byte(`[{"download_url":"a"},{"download_url":"b"}]`))
		case "/limited":
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message":"API rate limit exceeded"}`))
		default:
			_, _ = w.Write([]byte(`{"message":"Not Found"}`))
		}
	})
	got, err := GetRepo(srv.URL + "/ok")
	if err != nil || len(got) != 2 || got[0] != "a" {
		t.Fatalf("GetRepo = %v, %v", got, err)
	}
	if _, err := GetRepo(srv.URL + "/limited"); err == nil || !strings.Contains(err.Error(), "rate limit") {
		t.Fatalf("expected rate limit error, got %v", err)
	}
	if _, err := GetRepo(srv.URL + "/object"); err == nil {
		t.Fatal("expected decode error for non-array response")
	}
}

func TestBiliUpstreamErrorCode(t *testing.T) {
	srv := newServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":-404,"message":"啥都木有"}`))
	})
	orig := bilibiliViewAPI
	bilibiliViewAPI = srv.URL
	t.Cleanup(func() { bilibiliViewAPI = orig })
	if _, err := Bili("BV1xx411c7mD"); err == nil || !strings.Contains(err.Error(), "-404") {
		t.Fatalf("expected bilibili error code, got %v", err)
	}
}

func TestRedactURL(t *testing.T) {
	tests := map[string]string{
		"https://user:pass@example.com/a?key=secret#frag": "https://example.com/a",
		"http://example.com":                              "http://example.com",
		"::bad":                                           "<invalid URL>",
		"":                                                "<invalid URL>",
	}
	for in, want := range tests {
		if got := RedactURL(in); got != want {
			t.Errorf("RedactURL(%q) = %q, want %q", in, got, want)
		}
	}
}
