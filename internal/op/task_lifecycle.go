package op

import (
	"errors"
	"os"
	"sync"
	"time"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

// maxConcurrentPerAPI bounds simultaneous upstream generations per provider.
const maxConcurrentPerAPI = 3

// cacheSweepInterval is the minimum time between sweeps of expired results.
const cacheSweepInterval = time.Minute

// flight is one in-progress generation that identical requests wait on.
// task, result and err are written by the leader before done is closed.
type flight struct {
	done    chan struct{}
	waiters int
	task    model.Task
	result  GenerateResult
	err     error
}

type cachedResult struct {
	task   model.Task
	result GenerateResult
}

// taskCoordinator deduplicates identical in-flight work, caches successful
// results and limits concurrent generations per API. Its mutex is never held
// while generating.
type taskCoordinator struct {
	mu        sync.Mutex
	inflight  map[TaskQueueFilter]*flight
	cache     map[TaskQueueFilter]cachedResult
	limiters  map[string]chan struct{}
	limit     int
	lastSweep time.Time
	now       func() time.Time
}

func newTaskCoordinator(limit int) *taskCoordinator {
	return &taskCoordinator{
		inflight: make(map[TaskQueueFilter]*flight),
		cache:    make(map[TaskQueueFilter]cachedResult),
		limiters: make(map[string]chan struct{}),
		limit:    limit,
		now:      time.Now,
	}
}

var tasks = newTaskCoordinator(maxConcurrentPerAPI)

func cacheDuration(kind string) time.Duration {
	settings := database.SettingsStore.Get()
	minutes := settings.Text.CacheMinutes
	switch kind {
	case "img.gen":
		minutes = settings.Image.CacheMinutes
	case "web.img":
		minutes = settings.Web.CacheMinutes
	}
	return time.Duration(minutes) * time.Minute
}

func ExecuteCachedTask(task model.Task, filter TaskQueueFilter, skipDB bool, process func(*model.Task) (GenerateResult, error)) (model.Task, GenerateResult, error) {
	return tasks.execute(task, filter, skipDB, process)
}

func SaveTask(task model.Task, skipDB bool) {
	if !skipDB {
		util.ErrorPrinter(db.CreateTask(&task))
	}
}

// reuse turns a cached or shared successful result into this request's task.
func reuse(task model.Task, source model.Task) model.Task {
	id := task.UUID
	task = source
	task.UUID = id
	task.Time = time.Now()
	task.Temp = "Yes"
	return task
}

func (c *taskCoordinator) execute(task model.Task, filter TaskQueueFilter, skipDB bool, process func(*model.Task) (GenerateResult, error)) (model.Task, GenerateResult, error) {
	expire := cacheDuration(filter.Type)
	var staleImages []string

	c.mu.Lock()
	staleImages = c.sweepLocked()
	if entry, ok := c.cache[filter]; ok {
		if c.now().Sub(entry.task.Time) <= expire {
			c.mu.Unlock()
			removeImages(staleImages)
			reused := reuse(task, entry.task)
			SaveTask(reused, skipDB)
			return reused, entry.result, nil
		}
		delete(c.cache, filter)
		staleImages = append(staleImages, entry.task.UUID)
	}
	if f, ok := c.inflight[filter]; ok {
		f.waiters++
		c.mu.Unlock()
		removeImages(staleImages)
		<-f.done
		if f.err == nil && f.task.Status == "success" {
			reused := reuse(task, f.task)
			SaveTask(reused, skipDB)
			return reused, f.result, nil
		}
		failed := failTask(task, f.task.Return)
		failed.Temp = "No"
		SaveTask(failed, skipDB)
		err := f.err
		if err == nil {
			err = errors.New(f.task.Return)
		}
		return failed, GenerateResult{}, err
	}
	f := &flight{done: make(chan struct{})}
	c.inflight[filter] = f
	c.mu.Unlock()
	removeImages(staleImages)

	task.Temp = "No"
	return c.lead(f, task, filter, skipDB, process)
}

// lead runs the generation for a flight. Deferred cleanup guarantees that
// the concurrency permit is released and waiters are woken even if process
// panics.
func (c *taskCoordinator) lead(f *flight, task model.Task, filter TaskQueueFilter, skipDB bool, process func(*model.Task) (GenerateResult, error)) (model.Task, GenerateResult, error) {
	finished := false
	defer func() {
		if !finished {
			f.task = failTask(task, "Task aborted")
			f.err = errors.New("task aborted")
			c.complete(filter, f)
		}
	}()

	release := c.acquire(filter.API)
	result, err := func() (GenerateResult, error) {
		defer release()
		return process(&task)
	}()
	finishTask(&task, &result)

	f.task, f.result, f.err = task, result, err
	finished = true
	c.complete(filter, f)
	SaveTask(task, skipDB)
	return task, result, err
}

// complete publishes a flight's outcome, caching it on success, and wakes
// every waiter at once.
func (c *taskCoordinator) complete(filter TaskQueueFilter, f *flight) {
	c.mu.Lock()
	if c.inflight[filter] == f {
		delete(c.inflight, filter)
	}
	if f.err == nil && f.task.Status == "success" {
		c.cache[filter] = cachedResult{task: f.task, result: f.result}
	}
	c.mu.Unlock()
	close(f.done)
}

// acquire blocks until a generation slot for api is free and returns the
// function that releases it.
func (c *taskCoordinator) acquire(api string) func() {
	c.mu.Lock()
	sem, ok := c.limiters[api]
	if !ok {
		sem = make(chan struct{}, c.limit)
		c.limiters[api] = sem
	}
	c.mu.Unlock()
	sem <- struct{}{}
	var once sync.Once
	return func() { once.Do(func() { <-sem }) }
}

// sweepLocked drops expired cached results at most once per
// cacheSweepInterval and returns the image ids that are no longer served.
// The caller must hold c.mu.
func (c *taskCoordinator) sweepLocked() []string {
	now := c.now()
	if now.Sub(c.lastSweep) < cacheSweepInterval {
		return nil
	}
	c.lastSweep = now
	var stale []string
	for filter, entry := range c.cache {
		if now.Sub(entry.task.Time) > cacheDuration(filter.Type) {
			delete(c.cache, filter)
			stale = append(stale, entry.task.UUID)
		}
	}
	return stale
}

func removeImages(ids []string) {
	for _, id := range ids {
		if id != "" {
			_ = os.Remove(ImgPath + id + ".png")
		}
	}
}

// finishTask marks a nominally successful task as failed when it did not
// produce a valid image.
func finishTask(task *model.Task, result *GenerateResult) {
	if task.Status == "success" && task.Temp != "Yes" && !util.PngChecker(ImgPath+task.UUID+".png") {
		task.Status = "failed"
		task.Return = "Invalid Image File"
		result.URL = downloadURL(EmptyImageID)
	}
}
