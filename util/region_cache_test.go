package util

import (
	"fmt"
	"sync"
	"testing"
)

func TestRegionCacheBounded(t *testing.T) {
	c := &regionCache{m: make(map[string]string)}
	for i := 0; i < regionCacheLimit*2+5; i++ {
		c.set(fmt.Sprintf("10.0.%d.%d", i/256, i%256), "x")
		if c.len() > regionCacheLimit {
			t.Fatalf("cache grew to %d entries", c.len())
		}
	}
	c.set("1.1.1.1", "Earth")
	if v, ok := c.get("1.1.1.1"); !ok || v != "Earth" {
		t.Fatalf("get = %q, %v", v, ok)
	}
}

func TestRegionCacheConcurrentAccess(t *testing.T) {
	c := &regionCache{m: make(map[string]string)}
	var wg sync.WaitGroup
	for w := 0; w < 16; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				ip := fmt.Sprintf("192.168.%d.%d", w, i%256)
				c.set(ip, "r")
				c.get(ip)
			}
		}(w)
	}
	wg.Wait()
}
