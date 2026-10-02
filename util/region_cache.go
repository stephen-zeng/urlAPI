package util

import "sync"

// regionCacheLimit bounds the IP→region cache; when full it is cleared.
const regionCacheLimit = 1000

// regionCache is a small concurrency-safe IP→region cache.
type regionCache struct {
	mu sync.RWMutex
	m  map[string]string
}

var ipRegions = &regionCache{m: make(map[string]string)}

func (c *regionCache) get(ip string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	region, ok := c.m[ip]
	return region, ok
}

func (c *regionCache) set(ip, region string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.m) >= regionCacheLimit {
		c.m = make(map[string]string)
	}
	c.m[ip] = region
}

func (c *regionCache) len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.m)
}
