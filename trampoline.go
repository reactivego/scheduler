package scheduler

import (
	"sort"
	"time"
)

// Trampoline scheduler schedules a task to occur after the currently
// running task completes. A task scheduled on an empty trampoline
// will be dispatched sychronously and run immediately, while tasks
// scheduled by that task will be dispatched asynchronously and serial.
//
// Trampoline scheduler is not safe to use from multiple goroutines at
// the same time. It should be used purely for scheduling tasks from a
// single goroutine.
type Trampoline struct {
	*trampoline
	cancel chan struct{}
}

type trampoline struct {
	tasks []task
}

type task struct {
	at  time.Time
	run func()
}

// MakeTrampoline creates a new Trampoline scheduler instance.
func MakeTrampoline() *Trampoline {
	return &Trampoline{&trampoline{}, make(chan struct{})}
}

func (s *Trampoline) Add() *Trampoline {
	return &Trampoline{s.trampoline, make(chan struct{})}
}

// Len for sort.Sort support
func (s *Trampoline) Len() int {
	return len(s.tasks)
}

// Less for sort.Sort support
func (s *Trampoline) Less(i, j int) bool {
	return s.tasks[i].at.Before(s.tasks[j].at)
}

// Swap for sort.Sort support
func (s *Trampoline) Swap(i, j int) {
	s.tasks[i], s.tasks[j] = s.tasks[j], s.tasks[i]
}

// Now returns the current time according to the scheduler.
func (s *Trampoline) Now() time.Time {
	return time.Now()
}

// Schedule will dispatch the first task synchronously and any subsequent
// tasks asynchronously on a task queue. So when the first task eventually
// returns the queue of tasks is empty again.
func (s *Trampoline) Schedule(task func()) {
	s.ScheduleFuture(0, task)
}

// ScheduleRecursive will dispatch the first task synchronously
// and any subsequent tasks asynchronously on a task queue. So when
// the first task eventually returns the queue of tasks is empty again.
func (s *Trampoline) ScheduleRecursive(task func(self func())) {
	self := func() {
		s.ScheduleRecursive(task)
	}
	s.ScheduleFuture(0, func() {
		task(self)
	})
}

// ScheduleFuture will dispatch the first task synchronously and any subsequent
// tasks asynchronously on a task queue. So when the first task eventually
// returns the queue of tasks is empty again. The due parameter determines how
// far in the future the task will be scheduled.
func (s *Trampoline) ScheduleFuture(due time.Duration, taskfunc func()) {
	s.tasks = append(s.tasks, task{at: s.Now().Add(due), run: taskfunc})
	sort.Stable(s)
	if len(s.tasks) == 1 {
		for len(s.tasks) > 0 {
			select {
			case _, ok := <-s.cancel:
				if !ok {
					s.tasks = nil
					return // canceled
				}
			default:
			}
			task := s.tasks[0]

			now := time.Now()

			if now.Before(task.at) {
				timer := time.AfterFunc(task.at.Sub(now), func() {
					s.cancel <- struct{}{}
				})
				if _, ok := <-s.cancel; !ok {
					timer.Stop()
					s.tasks = nil
					return // canceled
				}
			}
			task.run()
			s.tasks = s.tasks[1:]
		}
	}
}

// ScheduleFutureRecursive will dispatch the first task synchronously
// and any subsequent tasks asynchronously on a task queue. When the first
// task eventually returns the queue of tasks is empty again.
// The task will be scheduled after the time period has passed.
// To schedule for the next period, the code should call the self function.
// Just returning from the task function will terminate a task.
func (s *Trampoline) ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration))) {
	self := func(due time.Duration) {
		s.ScheduleFutureRecursive(due, task)
	}
	s.ScheduleFuture(due, func() {
		task(self)
	})
}

// Cancel will remove all queued tasks from the scheduler. A running task is
// not affected by cancel and will continue until ist is finished.
func (s *Trampoline) Cancel() {
	if s.cancel != nil {
		close(s.cancel)
	}
}

// IsAsynchronous returns false when no task is currently scheduled.
func (s Trampoline) IsAsynchronous() bool {
	return len(s.tasks) > 0
}

// IsSerial returns false when no task is currently scheduled.
func (s Trampoline) IsSerial() bool {
	return len(s.tasks) > 0
}

// IsConcurrent returns false.
func (s Trampoline) IsConcurrent() bool {
	return false
}

func (s Trampoline) String() string {
	str := "Trampoline"
	if s.IsAsynchronous() {
		str += "{ Asynchronous:"
	} else {
		str += "{ Synchronous:"
	}
	if s.IsSerial() {
		str += "Serial }"
	} else {
		str += "Immediate }"
	}
	return str
}
