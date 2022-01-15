package scheduler

import (
	"fmt"
	"runtime"
	"sort"
	"time"
)

// futuretask

type futuretask struct {
	at     time.Time
	run    func()
	cancel chan struct{}
}

func (t *futuretask) Cancel() {
	if t.cancel != nil {
		close(t.cancel)
	}
}

// trampoline

type trampoline struct {
	gid     string
	tasks   []futuretask
	current *futuretask
}

// New creates and returns a serial (non-concurrent) scheduler that runs all
// tasks on a single goroutine. The returned scheduler is returned as a Scheduler
// interface. Tasks scheduled will be dispatched asynchronously because they are
// added to a serial queue. When the Wait method is called all tasks scheduled
// on the serial queue will be performed in the same order in which they were added
// to the queue.
//
// The returned scheduler is not safe to be shared by multiple goroutines
// concurrently. It should be used purely from a single goroutine to schedule
// tasks to run sequentially.
func New() Scheduler {
	return &trampoline{gid: Gid()}
}

// MakeTrampoline is deprecated, use New instead
var MakeTrampoline = New

func (s *trampoline) Len() int {
	return len(s.tasks)
}

func (s *trampoline) Less(i, j int) bool {
	return s.tasks[i].at.Before(s.tasks[j].at)
}

func (s *trampoline) Swap(i, j int) {
	s.tasks[i], s.tasks[j] = s.tasks[j], s.tasks[i]
}

func (s *trampoline) Now() time.Time {
	return time.Now()
}

func (s *trampoline) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (s *trampoline) Schedule(task func()) Runner {
	t := futuretask{at: time.Now(), run: task, cancel: make(chan struct{})}
	s.tasks = append(s.tasks, t)
	sort.Stable(s)
	return &t
}

func (s *trampoline) ScheduleRecursive(task func(again func())) Runner {
	t := futuretask{cancel: make(chan struct{})}
	again := func() {
		t.at = time.Now()
		s.tasks = append(s.tasks, t)
		sort.Stable(s)
	}
	t.run = func() {
		task(again)
	}
	again()
	return &t
}

func (s *trampoline) ScheduleLoop(from int, task func(index int, again func(next int))) Runner {
	t := futuretask{cancel: make(chan struct{})}
	var run func(index int) func()
	again := func(index int) {
		t.at = time.Now()
		t.run = run(index)
		s.tasks = append(s.tasks, t)
		sort.Stable(s)
	}
	run = func(index int) func() {
		return func() { task(index, again) }
	}
	again(from)
	return &t
}

func (s *trampoline) ScheduleFuture(due time.Duration, task func()) Runner {
	t := futuretask{at: time.Now().Add(due), run: task, cancel: make(chan struct{})}
	s.tasks = append(s.tasks, t)
	sort.Stable(s)
	return &t
}

func (s *trampoline) ScheduleFutureRecursive(due time.Duration, task func(again func(time.Duration))) Runner {
	t := futuretask{cancel: make(chan struct{})}
	again := func(due time.Duration) {
		t.at = time.Now().Add(due)
		s.tasks = append(s.tasks, t)
		sort.Stable(s)
	}
	t.run = func() {
		task(again)
	}
	again(due)
	return &t
}

func (s *trampoline) Wait() {
	for s.RunTask() {
	}
}

func (s *trampoline) Gosched() {
	if len(s.gid) > 0 && s.gid == Gid() {
		if s.RunTask() {
			return
		}
	}
	runtime.Gosched()
}

func (s *trampoline) RunTask() bool {
	if len(s.tasks) == 0 {
		return false
	}
	s.current = &s.tasks[0]
	s.tasks = s.tasks[1:]
	if time.Until(s.current.at) < 999*time.Millisecond {
		s.ShortWaitAndRun(s.current)
	} else {
		s.LongWaitAndRun(s.current)
	}
	s.current = nil
	return true
}

func (s *trampoline) ShortWaitAndRun(task *futuretask) {
	for time.Now().Before(task.at) {
		select {
		case <-task.cancel:
			return
		default:
			runtime.Gosched()
		}
	}
	select {
	case <-task.cancel:
		return
	default:
		task.run()
	}
}

func (s *trampoline) LongWaitAndRun(task *futuretask) {
	due := time.Until(task.at)
	if due > 0 {
		deadline := time.NewTimer(due)
		select {
		case <-task.cancel:
			deadline.Stop()
			return
		case <-deadline.C:
			task.run()
			return
		}
	}
	select {
	case <-task.cancel:
		return
	default:
		task.run()
	}
}

func (s *trampoline) IsConcurrent() bool {
	return false
}

func (s *trampoline) Count() int {
	if s.current == nil {
		return len(s.tasks)
	} else {
		return len(s.tasks) + 1
	}
}

func (s trampoline) String() string {
	at := make([]string, len(s.tasks))
	for i := range s.tasks {
		at[i] = s.tasks[i].at.Format("15:04:05")
	}
	return fmt.Sprintf("Trampoline{ gid = %s, tasks = %d, at = %v }", s.gid, len(s.tasks), at)
}
