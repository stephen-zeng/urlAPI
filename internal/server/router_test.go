package server

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestServeStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	srv := &http.Server{Addr: "127.0.0.1:0", Handler: http.NotFoundHandler()}
	done := make(chan error, 1)
	go func() { done <- serve(ctx, srv) }()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("serve returned %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serve did not stop after cancellation")
	}
}

func TestServeReturnsListenError(t *testing.T) {
	srv := &http.Server{Addr: "invalid-address:-1"}
	if err := serve(context.Background(), srv); err == nil {
		t.Fatal("expected listen error")
	}
}
