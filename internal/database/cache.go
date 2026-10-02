package database

import (
	"sync"
	"urlAPI/internal/model"
)

// SessionCache is a concurrency-safe in-memory index of login sessions keyed
// by token. It mirrors the sessions table so authentication avoids a query.
type SessionCache struct {
	mu sync.RWMutex
	m  map[string]model.Session
}

func newSessionCache() *SessionCache {
	return &SessionCache{m: make(map[string]model.Session)}
}

func (c *SessionCache) Get(token string) (model.Session, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	session, ok := c.m[token]
	return session, ok
}

func (c *SessionCache) Set(session model.Session) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[session.Token] = session
}

func (c *SessionCache) Delete(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, token)
}

func (c *SessionCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.m)
}

// Reset replaces the whole cache with the given sessions.
func (c *SessionCache) Reset(sessions []model.Session) {
	m := make(map[string]model.Session, len(sessions))
	for _, session := range sessions {
		m[session.Token] = session
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m = m
}

// RepoCache is a concurrency-safe in-memory index of random-image repository
// contents keyed by API and repository info. Stored slices are never mutated
// after insertion, so callers may read a returned slice without locking.
type RepoCache struct {
	mu sync.RWMutex
	m  map[string][]string
}

func newRepoCache() *RepoCache {
	return &RepoCache{m: make(map[string][]string)}
}

func repoKey(api, info string) string {
	return api + ";" + info
}

func (c *RepoCache) Get(api, info string) ([]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	content, ok := c.m[repoKey(api, info)]
	return content, ok
}

func (c *RepoCache) Set(api, info string, content []string) {
	stored := append([]string(nil), content...)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[repoKey(api, info)] = stored
}

func (c *RepoCache) Delete(api, info string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, repoKey(api, info))
}

func (c *RepoCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.m)
}
