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
type Trampoline struct{ tasks []Task }

type Task struct {
	at  time.Time
	run func()
}

// Len for sort.Sort support
func (t *Trampoline) Len() int {
	return len(t.tasks)
}

// Less for sort.Sort support
func (t *Trampoline) Less(i, j int) bool {
	return t.tasks[i].at.Before(t.tasks[j].at)
}

// Swap for sort.Sort support
func (t *Trampoline) Swap(i, j int) {
	t.tasks[i], t.tasks[j] = t.tasks[j], t.tasks[i]
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
func (s *Trampoline) ScheduleFuture(due time.Duration, task func()) {
	s.tasks = append(s.tasks, Task{at: s.Now().Add(due), run: task})
	sort.Stable(s)
	if len(s.tasks) == 1 {
		for len(s.tasks) > 0 {
			// Wait until the at time has arrived....
			task := s.tasks[0]
			now := time.Now()
			if now.Before(task.at) {
				time.Sleep(task.at.Sub(now))
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

// Cancel will remove all queued tasks from the scheduler and stop a wait for
// the next due time when active. A running task is not affected by cancel and
// will continue until it is finished.
func (s *Trampoline) Cancel() {
	s.tasks = nil
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
