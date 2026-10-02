package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	// maxAPIResponseBytes bounds JSON and HTML bodies read from upstreams.
	maxAPIResponseBytes = 8 << 20
	// maxDownloadBytes bounds downloaded images.
	maxDownloadBytes = 32 << 20
	// maxErrorDetail bounds the upstream message embedded in errors.
	maxErrorDetail = 200
)

// ErrResponseTooLarge is returned when an upstream body exceeds its limit.
var ErrResponseTooLarge = errors.New("upstream response too large")

// UpstreamError describes an unexpected (non-2xx) upstream HTTP response.
// Detail holds a short, credential-redacted message extracted from the body.
type UpstreamError struct {
	StatusCode int
	Status     string
	Detail     string
}

func (e *UpstreamError) Error() string {
	if e.Detail == "" {
		return "upstream returned " + e.Status
	}
	return "upstream returned " + e.Status + ": " + e.Detail
}

// upstreamRequest describes one outbound call.
type upstreamRequest struct {
	Method  string
	URL     string
	Body    []byte
	Headers map[string]string
	// Secrets are redacted from any error text derived from the response.
	Secrets []string
	// Limit caps the response body size; maxAPIResponseBytes when zero.
	Limit int64
}

// doUpstream performs the request with the shared client and returns the
// body of a 2xx response. Every failure yields a non-nil error that never
// contains request headers, query strings or the configured secrets.
func doUpstream(ctx context.Context, r upstreamRequest) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	method := r.Method
	if method == "" {
		method = http.MethodGet
	}
	var body io.Reader
	if r.Body != nil {
		body = bytes.NewReader(r.Body)
	}
	req, err := http.NewRequestWithContext(ctx, method, r.URL, body)
	if err != nil {
		return nil, errors.Errorf("build request for %s: invalid URL", RedactURL(r.URL))
	}
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}
	resp, err := GlobalHTTPClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(sanitizeTransportError(err, r.Secrets))
	}
	defer resp.Body.Close()

	limit := r.Limit
	if limit <= 0 {
		limit = maxAPIResponseBytes
	}
	data, readErr := readLimited(resp.Body, limit)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.WithStack(&UpstreamError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Detail:     extractErrorDetail(data, r.Secrets),
		})
	}
	if readErr != nil {
		return nil, errors.Wrapf(readErr, "read response from %s", RedactURL(r.URL))
	}
	return data, nil
}

// doUpstreamJSON performs the request and decodes a 2xx JSON response.
func doUpstreamJSON(ctx context.Context, r upstreamRequest, out any) error {
	data, err := doUpstream(ctx, r)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, out); err != nil {
		return errors.Wrapf(err, "decode response from %s", RedactURL(r.URL))
	}
	return nil
}

func readLimited(r io.Reader, limit int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return data, err
	}
	if int64(len(data)) > limit {
		return data[:limit], ErrResponseTooLarge
	}
	return data, nil
}

// RedactURL removes user info, query and fragment from a URL so it can be
// safely included in logs and errors.
func RedactURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return "<invalid URL>"
	}
	u.User = nil
	u.RawQuery = ""
	u.ForceQuery = false
	u.Fragment = ""
	u.RawFragment = ""
	return u.String()
}

// sanitizeTransportError rewrites *url.Error so the full request URL (which
// may carry an API key in its query) is not exposed.
func sanitizeTransportError(err error, secrets []string) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		msg := fmt.Sprintf("%s %s: %s", urlErr.Op, RedactURL(urlErr.URL), redactSecrets(urlErr.Err.Error(), secrets))
		if urlErr.Timeout() {
			return &timeoutError{msg: msg}
		}
		return errors.New(msg)
	}
	return errors.New(redactSecrets(err.Error(), secrets))
}

type timeoutError struct{ msg string }

func (e *timeoutError) Error() string { return e.msg }
func (e *timeoutError) Timeout() bool { return true }

func redactSecrets(s string, secrets []string) string {
	for _, secret := range secrets {
		if len(secret) >= 4 {
			s = strings.ReplaceAll(s, secret, "[REDACTED]")
		}
	}
	return s
}

// extractErrorDetail pulls a human readable message out of common upstream
// error payloads ({"error":{"message":...}}, {"message":...}, ...).
func extractErrorDetail(data []byte, secrets []string) string {
	var payload struct {
		Error   json.RawMessage `json:"error"`
		Message string          `json:"message"`
		Code    any             `json:"code"`
	}
	var detail string
	if json.Unmarshal(data, &payload) == nil {
		var nested struct {
			Message string `json:"message"`
		}
		var plain string
		switch {
		case len(payload.Error) > 0 && json.Unmarshal(payload.Error, &nested) == nil && nested.Message != "":
			detail = nested.Message
		case len(payload.Error) > 0 && json.Unmarshal(payload.Error, &plain) == nil && plain != "":
			detail = plain
		case payload.Message != "":
			detail = payload.Message
		}
	}
	detail = strings.Join(strings.Fields(redactSecrets(detail, secrets)), " ")
	if r := []rune(detail); len(r) > maxErrorDetail {
		detail = string(r[:maxErrorDetail]) + "…"
	}
	return detail
}
