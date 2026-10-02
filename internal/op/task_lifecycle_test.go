package op

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

// setupCoordinator isolates the global coordinator, image directory and
// settings for one test.
func setupCoordinator(t *testing.T, cacheMinutes int) *taskCoordinator {
	t.Helper()
	origTasks, origPath, origSettings := tasks, ImgPath, database.SettingsStore.Get()
	t.Cleanup(func() {
		tasks, ImgPath = origTasks, origPath
		database.SettingsStore.Replace(origSettings)
	})
	ImgPath = t.TempDir() + "/"
	settings := util.AppSettings{}
	settings.Text.CacheMinutes = cacheMinutes
	settings.Image.CacheMinutes = cacheMinutes
	settings.Web.CacheMinutes = cacheMinutes
	database.SettingsStore.Replace(settings)
	tasks = newTaskCoordinator(maxConcurrentPerAPI)
	return tasks
}

func writePNG(t *testing.T, id string) {
	t.Helper()
	f, err := os.Create(ImgPath + id + ".png")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Error(err)
	}
}

// succeed returns a process function that writes a valid image.
func succeed(t *testing.T, calls *atomic.Int32, gate <-chan struct{}) func(*model.Task) (GenerateResult, error) {
	return func(task *model.Task) (GenerateResult, error) {
		calls.Add(1)
		if gate != nil {
			<-gate
		}
		writePNG(t, task.UUID)
		result := GenerateResult{URL: "/download?img=" + task.UUID}
		return result, setTaskResult(task, result)
	}
}

func newTask(i int) model.Task {
	return model.Task{UUID: fmt.Sprintf("uuid-%d", i), Time: time.Now()}
}

// waitFor polls cond until it holds or the deadline passes.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(time.Millisecond)
	}
}

func waiters(c *taskCoordinator, filter TaskQueueFilter) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if f, ok := c.inflight[filter]; ok {
		return f.waiters
	}
	return -1
}

func TestDuplicateWorkIsCoalesced(t *testing.T) {
	c := setupCoordinator(t, 10)
	filter := TaskQueueFilter{Type: "txt.gen", Target: "hello", API: "openai"}
	var calls atomic.Int32
	gate := make(chan struct{})
	process := succeed(t, &calls, gate)

	const n = 8
	type out struct {
		task model.Task
		res  GenerateResult
		err  error
	}
	results := make(chan out, n)
	go func() {
		task, res, err := ExecuteCachedTask(newTask(0), filter, true, process)
		results <- out{task, res, err}
	}()
	waitFor(t, "leader to start", func() bool { return calls.Load() == 1 })
	for i := 1; i < n; i++ {
		go func(i int) {
			task, res, err := ExecuteCachedTask(newTask(i), filter, true, process)
			results <- out{task, res, err}
		}(i)
	}
	waitFor(t, "waiters to block", func() bool { return waiters(c, filter) == n-1 })
	close(gate)

	reused := 0
	for i := 0; i < n; i++ {
		r := <-results
		if r.err != nil || r.task.Status != "success" || r.res.URL != "/download?img=uuid-0" {
			t.Fatalf("result %d: %+v %v", i, r, r.err)
		}
		if r.task.Temp == "Yes" {
			reused++
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("process ran %d times, want 1", calls.Load())
	}
	if reused != n-1 {
		t.Fatalf("%d results reused, want %d", reused, n-1)
	}
}

func TestProviderConcurrencyLimit(t *testing.T) {
	setupCoordinator(t, 10)
	var active, maxActive atomic.Int32
	gate := make(chan struct{})
	process := func(task *model.Task) (GenerateResult, error) {
		cur := active.Add(1)
		for {
			prev := maxActive.Load()
			if cur <= prev || maxActive.CompareAndSwap(prev, cur) {
				break
			}
		}
		<-gate
		active.Add(-1)
		writePNG(t, task.UUID)
		return GenerateResult{}, setTaskResult(task, GenerateResult{})
	}

	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			filter := TaskQueueFilter{Type: "img.gen", Target: fmt.Sprint(i), API: "alibaba"}
			if _, _, err := ExecuteCachedTask(newTask(i), filter, true, process); err != nil {
				t.Error(err)
			}
		}(i)
	}
	// Different APIs have independent limits.
	other := make(chan error, 1)
	go func() {
		_, _, err := ExecuteCachedTask(newTask(99), TaskQueueFilter{Type: "img.gen", Target: "x", API: "openai"}, true,
			func(task *model.Task) (GenerateResult, error) {
				writePNG(t, task.UUID)
				return GenerateResult{}, setTaskResult(task, GenerateResult{})
			})
		other <- err
	}()

	waitFor(t, "slots to fill", func() bool { return active.Load() == maxConcurrentPerAPI })
	select {
	case err := <-other:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("a different API was blocked by the alibaba limit")
	}
	time.Sleep(20 * time.Millisecond)
	if got := active.Load(); got != maxConcurrentPerAPI {
		t.Fatalf("active = %d, want %d", got, maxConcurrentPerAPI)
	}
	close(gate)
	wg.Wait()
	if got := maxActive.Load(); got != maxConcurrentPerAPI {
		t.Fatalf("max concurrent = %d, want %d", got, maxConcurrentPerAPI)
	}
}

func TestFailureWakesWaitersAndReleasesSlot(t *testing.T) {
	c := setupCoordinator(t, 10)
	filter := TaskQueueFilter{Type: "txt.gen", Target: "boom", API: "deepseek"}
	var calls atomic.Int32
	gate := make(chan struct{})
	fail := func(task *model.Task) (GenerateResult, error) {
		calls.Add(1)
		<-gate
		task.Status = "failed"
		task.Return = "upstream exploded"
		return GenerateResult{}, errors.New("upstream exploded")
	}

	const n = 5
	errs := make(chan error, n)
	go func() { _, _, err := ExecuteCachedTask(newTask(0), filter, true, fail); errs <- err }()
	waitFor(t, "leader to start", func() bool { return calls.Load() == 1 })
	for i := 1; i < n; i++ {
		go func(i int) {
			task, _, err := ExecuteCachedTask(newTask(i), filter, true, fail)
			if task.Status != "failed" || task.Return != "upstream exploded" {
				t.Errorf("waiter task = %+v", task)
			}
			errs <- err
		}(i)
	}
	waitFor(t, "waiters to block", func() bool { return waiters(c, filter) == n-1 })
	close(gate)
	for i := 0; i < n; i++ {
		select {
		case err := <-errs:
			if err == nil {
				t.Fatal("expected failure to propagate")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("waiter was not woken after failure")
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("process ran %d times, want 1", calls.Load())
	}

	// Failures are not cached and must not hold slots: maxConcurrentPerAPI
	// new generations for the same API can all run at once.
	var active atomic.Int32
	release := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < maxConcurrentPerAPI; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			f := filter
			if i > 0 {
				f.Target = fmt.Sprint(i)
			}
			_, _, _ = ExecuteCachedTask(newTask(100+i), f, true, func(task *model.Task) (GenerateResult, error) {
				active.Add(1)
				<-release
				return GenerateResult{}, errors.New("again")
			})
		}(i)
	}
	waitFor(t, "all slots to be usable after failures", func() bool { return active.Load() == maxConcurrentPerAPI })
	close(release)
	wg.Wait()
}

func TestPanicReleasesSlotAndWakesWaiters(t *testing.T) {
	c := setupCoordinator(t, 10)
	c.limit = 1
	filter := TaskQueueFilter{Type: "web.img", Target: "p", API: "github.com"}
	gate := make(chan struct{})
	var calls atomic.Int32

	leaderDone := make(chan any, 1)
	go func() {
		defer func() { leaderDone <- recover() }()
		_, _, _ = ExecuteCachedTask(newTask(0), filter, true, func(*model.Task) (GenerateResult, error) {
			calls.Add(1)
			<-gate
			panic("kaboom")
		})
	}()
	waitFor(t, "leader to start", func() bool { return calls.Load() == 1 })
	waiterErr := make(chan error, 1)
	go func() {
		_, _, err := ExecuteCachedTask(newTask(1), filter, true, nil)
		waiterErr <- err
	}()
	waitFor(t, "waiter to block", func() bool { return waiters(c, filter) == 1 })
	close(gate)
	if r := <-leaderDone; r == nil {
		t.Fatal("panic was swallowed")
	}
	if err := <-waiterErr; err == nil {
		t.Fatal("waiter should observe the aborted task")
	}

	// With a limit of 1, this only completes if the panicking leader
	// released its slot.
	var calls2 atomic.Int32
	done := make(chan struct{})
	go func() {
		_, _, _ = ExecuteCachedTask(newTask(2), filter, true, succeed(t, &calls2, nil))
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("slot leaked after panic")
	}
}

func TestCachedResultsAreReusedUntilExpiry(t *testing.T) {
	c := setupCoordinator(t, 10)
	now := time.Now()
	c.now = func() time.Time { return now }
	filter := TaskQueueFilter{Type: "img.gen", Target: "cat", API: "openai"}
	var calls atomic.Int32
	process := succeed(t, &calls, nil)

	first, res1, err := ExecuteCachedTask(newTask(0), filter, true, process)
	if err != nil || first.Temp != "No" {
		t.Fatalf("first run: %+v %v", first, err)
	}
	second, res2, err := ExecuteCachedTask(newTask(1), filter, true, process)
	if err != nil || second.Temp != "Yes" || second.UUID != "uuid-1" || res2 != res1 {
		t.Fatalf("cached run: %+v %+v %v", second, res2, err)
	}
	if calls.Load() != 1 {
		t.Fatalf("process ran %d times, want 1", calls.Load())
	}

	// Advance past the 10 minute cache window: the result is regenerated and
	// the stale image removed.
	now = now.Add(11 * time.Minute)
	third, _, err := ExecuteCachedTask(newTask(2), filter, true, process)
	if err != nil || third.Temp != "No" || calls.Load() != 2 {
		t.Fatalf("expired run: %+v %v calls=%d", third, err, calls.Load())
	}
	if _, err := os.Stat(ImgPath + "uuid-0.png"); !os.IsNotExist(err) {
		t.Fatalf("stale image not removed: %v", err)
	}
}

func TestFailedResultsAreNotCached(t *testing.T) {
	setupCoordinator(t, 10)
	filter := TaskQueueFilter{Type: "txt.gen", Target: "x", API: "openai"}
	var calls atomic.Int32
	noImage := func(task *model.Task) (GenerateResult, error) {
		calls.Add(1)
		// Reports success but writes no image: finishTask must reject it.
		return GenerateResult{}, setTaskResult(task, GenerateResult{})
	}
	task, res, _ := ExecuteCachedTask(newTask(0), filter, true, noImage)
	if task.Status != "failed" || res.URL != "/download?img=empty" {
		t.Fatalf("invalid image accepted: %+v %+v", task, res)
	}
	_, _, _ = ExecuteCachedTask(newTask(1), filter, true, noImage)
	if calls.Load() != 2 {
		t.Fatalf("failed result was cached: calls=%d", calls.Load())
	}
}

func TestSweepBoundsCache(t *testing.T) {
	c := setupCoordinator(t, 1)
	now := time.Now()
	c.now = func() time.Time { return now }
	var calls atomic.Int32
	for i := 0; i < 50; i++ {
		filter := TaskQueueFilter{Type: "txt.gen", Target: fmt.Sprint(i), API: "openai"}
		task := newTask(i)
		task.Time = now
		if _, _, err := ExecuteCachedTask(task, filter, true, succeed(t, &calls, nil)); err != nil {
			t.Fatal(err)
		}
	}
	now = now.Add(2 * time.Hour)
	_, _, _ = ExecuteCachedTask(newTask(1000), TaskQueueFilter{Type: "txt.gen", Target: "new", API: "openai"}, true, succeed(t, &calls, nil))
	c.mu.Lock()
	n := len(c.cache)
	c.mu.Unlock()
	if n != 1 {
		t.Fatalf("cache holds %d entries after sweep, want 1", n)
	}
}
