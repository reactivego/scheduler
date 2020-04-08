package scheduler

import (
	"sort"
	"time"
)

// Trampoline is a serial (non-concurrent) scheduler that runs all tasks on 
// a single goroutine. The first task scheduled on an empty trampoline
// scheduler will run immediately and the Schedule function will return only
// once the task has finished. However, the tasks scheduled by that initial
// task will be dispatched asynchronously because they are added to a serial
// queue. Now when the first task is finished, before returning to the user
// all tasks scheduled on the serial queue will be performed in dispatch order.
// 
// The Trampoline scheduler is not safe to use from multiple goroutines at the
// same time. It should be used purely for scheduling tasks from a single
// goroutine.
var Trampoline = MakeTrampoline()

type trampoline struct {
	*trampo
	cancel chan struct{}
}

type trampo struct {
	tasks []task
}

type task struct {
	at  time.Time
	run func()
}

// MakeTrampoline creates and returns a new serial (non-concurrent) scheduler
// instance. The returned instance implements the Scheduler interface.
func MakeTrampoline() *trampoline {
	return &trampoline{&trampo{}, make(chan struct{})}
}

func (s *trampoline) Add() *trampoline {
	return &trampoline{s.trampo, make(chan struct{})}
}

// Len for sort.Sort support
func (s *trampoline) Len() int {
	return len(s.tasks)
}

// Less for sort.Sort support
func (s *trampoline) Less(i, j int) bool {
	return s.tasks[i].at.Before(s.tasks[j].at)
}

// Swap for sort.Sort support
func (s *trampoline) Swap(i, j int) {
	s.tasks[i], s.tasks[j] = s.tasks[j], s.tasks[i]
}

// Now returns the current time according to the scheduler.
func (s *trampoline) Now() time.Time {
	return time.Now()
}

// Schedule will dispatch the first task synchronously and any subsequent
// tasks asynchronously on a task queue. So when the first task eventually
// returns the queue of tasks is empty again.
func (s *trampoline) Schedule(task func()) {
	s.ScheduleFuture(0, task)
}

// ScheduleRecursive will dispatch the first task synchronously
// and any subsequent tasks asynchronously on a task queue. So when
// the first task eventually returns the queue of tasks is empty again.
func (s *trampoline) ScheduleRecursive(task func(self func())) {
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
func (s *trampoline) ScheduleFuture(due time.Duration, taskfunc func()) {
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
func (s *trampoline) ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration))) {
	self := func(due time.Duration) {
		s.ScheduleFutureRecursive(due, task)
	}
	s.ScheduleFuture(due, func() {
		task(self)
	})
}

// Cancel will remove all queued tasks from the scheduler. A running task is
// not affected by cancel and will continue until ist is finished.
func (s *trampoline) Cancel() {
	if s.cancel != nil {
		close(s.cancel)
	}
}

// IsAsynchronous returns false when no task is currently scheduled.
func (s trampoline) IsAsynchronous() bool {
	return len(s.tasks) > 0
}

func (s trampoline) String() string {
	str := "Trampoline"
	if s.IsAsynchronous() {
		str += "{ Asynchronous:Serial }"
	} else {
		str += "{ Synchronous:Immdediate }"
	}
	return str
}
