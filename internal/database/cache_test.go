package database

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"urlAPI/internal/model"
)

func TestSessionCacheBasic(t *testing.T) {
	c := newSessionCache()
	if _, ok := c.Get("missing"); ok {
		t.Fatal("unexpected hit on empty cache")
	}
	s := model.Session{Token: "tok", Expire: time.Now().Add(time.Hour), Term: true}
	c.Set(s)
	got, ok := c.Get("tok")
	if !ok || got.Token != "tok" || !got.Term {
		t.Fatalf("Get = %+v, %v", got, ok)
	}
	c.Delete("tok")
	if _, ok := c.Get("tok"); ok {
		t.Fatal("session still present after Delete")
	}
	c.Reset([]model.Session{{Token: "a"}, {Token: "b"}})
	if c.Len() != 2 {
		t.Fatalf("Len = %d, want 2", c.Len())
	}
	c.Reset(nil)
	if c.Len() != 0 {
		t.Fatalf("Len = %d after Reset(nil), want 0", c.Len())
	}
}

func TestRepoCacheCopiesInput(t *testing.T) {
	c := newRepoCache()
	content := []string{"a", "b"}
	c.Set("github", "owner/repo", content)
	content[0] = "mutated"
	got, ok := c.Get("github", "owner/repo")
	if !ok || got[0] != "a" {
		t.Fatalf("cache shares caller slice: %v", got)
	}
	c.Delete("github", "owner/repo")
	if _, ok := c.Get("github", "owner/repo"); ok {
		t.Fatal("repo still present after Delete")
	}
}

// Run with -race: concurrent readers and writers must not race.
func TestCachesConcurrentAccess(t *testing.T) {
	sessions := newSessionCache()
	repos := newRepoCache()
	const workers = 16
	const iterations = 200
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				token := fmt.Sprintf("t%d-%d", w, i%10)
				sessions.Set(model.Session{Token: token})
				sessions.Get(token)
				repos.Set("github", token, []string{token})
				if got, ok := repos.Get("github", token); ok && len(got) > 0 {
					_ = got[0]
				}
				if i%7 == 0 {
					sessions.Delete(token)
					repos.Delete("github", token)
				}
				if i%50 == 0 {
					sessions.Reset(nil)
				}
				sessions.Len()
				repos.Len()
			}
		}(w)
	}
	wg.Wait()
}
