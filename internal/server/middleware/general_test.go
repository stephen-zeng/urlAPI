package middleware

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func resetFrequency(t *testing.T) {
	t.Helper()
	IPFrequency.mu.Lock()
	IPFrequency.IPFrequency = make(map[FrequencyFilter]FrequencyData)
	IPFrequency.mu.Unlock()
}

func TestCheckFrequencyLimitsBursts(t *testing.T) {
	resetFrequency(t)
	now := time.Now()
	var last General
	for i := 0; i < 12; i++ {
		last = General{Type: "txt", IP: "1.2.3.4", Time: now}
		checkFrequency(&last)
	}
	if !last.Unsafe {
		t.Fatal("expected burst of requests to be flagged")
	}
	other := General{Type: "txt", IP: "5.6.7.8", Time: now}
	checkFrequency(&other)
	if other.Unsafe {
		t.Fatal("unrelated IP must not be limited")
	}
}

func TestCheckFrequencyPrunesExpiredEntries(t *testing.T) {
	resetFrequency(t)
	old := time.Now().Add(-time.Minute)
	IPFrequency.mu.Lock()
	for i := 0; i < frequencyPruneThreshold; i++ {
		IPFrequency.IPFrequency[FrequencyFilter{Type: "img", IP: fmt.Sprint(i)}] = FrequencyData{Counter: 1, Time: old}
	}
	IPFrequency.mu.Unlock()

	g := General{Type: "img", IP: "fresh", Time: time.Now()}
	checkFrequency(&g)

	IPFrequency.mu.Lock()
	n := len(IPFrequency.IPFrequency)
	IPFrequency.mu.Unlock()
	if n != 1 {
		t.Fatalf("tracker holds %d entries after pruning, want 1", n)
	}
}

func TestCheckFrequencyConcurrent(t *testing.T) {
	resetFrequency(t)
	var wg sync.WaitGroup
	for w := 0; w < 16; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				g := General{Type: "web", IP: fmt.Sprintf("10.0.0.%d", w), Time: time.Now()}
				checkFrequency(&g)
			}
		}(w)
	}
	wg.Wait()
}
