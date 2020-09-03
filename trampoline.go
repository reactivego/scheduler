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

// MakeTrampoline creates and returns a non-concurrent scheduler that runs
// all tasks on a single goroutine. The returned instance implements the
// Scheduler interface. Tasks scheduled will be dispatched asynchronously
// because they are added to a serial queue. Now when the Wait method is called
// all tasks scheduled on the serial queue will be performed in dispatch order.
//
// The Trampoline scheduler is not safe to use from multiple goroutines at the
// concurrently. It should be used purely for scheduling tasks from a single
// goroutine.
func MakeTrampoline() *trampoline {
	return &trampoline{gid: Gid()}
}

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
	t := futuretask{time.Now(), task, make(chan struct{})}
	s.tasks = append(s.tasks, t)
	sort.Stable(s)
	return &t
}

func (s *trampoline) ScheduleRecursive(task func(self func())) Runner {
	t := futuretask{cancel: make(chan struct{})}
	self := func() {
		t.at = time.Now()
		s.tasks = append(s.tasks, t)
		sort.Stable(s)
	}
	t.run = func() {
		task(self)
	}
	self()
	return &t
}

func (s *trampoline) ScheduleFuture(due time.Duration, task func()) Runner {
	t := futuretask{time.Now().Add(due), task, make(chan struct{})}
	s.tasks = append(s.tasks, t)
	sort.Stable(s)
	return &t
}

func (s *trampoline) ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration))) Runner {
	t := futuretask{cancel: make(chan struct{})}
	self := func(due time.Duration) {
		t.at = time.Now().Add(due)
		s.tasks = append(s.tasks, t)
		sort.Stable(s)
	}
	t.run = func() {
		task(self)
	}
	self(due)
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
	return fmt.Sprintf("Trampoline{ gid = %s, tasks = %d }", s.gid, len(s.tasks))
}
