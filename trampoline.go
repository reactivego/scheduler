package scheduler

import (
	"sync/atomic"
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
type Trampoline struct{ tasks []func() }

func (s *Trampoline) Now() time.Time {
	return time.Now()
}

// Schedule will dispatch the first task synchronously and any subsequent
// tasks asynchronously on a task queue. So when the first task eventually
// returns the queue of tasks is empty again.
func (s *Trampoline) Schedule(task func()) {
	s.tasks = append(s.tasks, task)
	if len(s.tasks) == 1 {
		for len(s.tasks) > 0 {
			s.tasks[0]()
			s.tasks = s.tasks[1:]
		}
	}
}

// ScheduleRecursive will dispatch the first task synchronously
// and any subsequent tasks asynchronously on a task queue. So when
// the first task eventually returns the queue of tasks is empty again.
func (s *Trampoline) ScheduleRecursive(task func(self func())) {
	self := func() {
		s.ScheduleRecursive(task)
	}
	s.Schedule(func() {
		task(self)
	})
}

// ScheduleFutureRecursive will dispatch the first task synchronously
// and any subsequent tasks asynchronously on a task queue. When the first
// task eventually returns the queue of tasks is empty again.
// The task will be scheduled after the time period has passed.
// To schedule for the next period, the code should call the self function.
// Just returning from the task function will terminate a task.
func (s *Trampoline) ScheduleFutureRecursive(timeout time.Duration, task func(self func(time.Duration))) {
	self := func(timeout time.Duration) {
		s.ScheduleFutureRecursive(timeout, task)
	}
	s.Schedule(func() {
		time.Sleep(timeout)
		task(self)
	})
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

// Wait will run the tasks from the queue by executing every task one by one.
// It will return when the function registered via onCancel() is called or
// when there are no more tasks remaining. Note, the currently running task
// may append additional tasks to the queue to run later.
func (s *Trampoline) Wait(onCancel func(func())) {
	active := int32(1)
	onCancel(func() {
		atomic.StoreInt32(&active, 0)
	})
	for atomic.LoadInt32(&active) != 0 && len(s.tasks) > 0 {
		s.tasks[0]()
		s.tasks = s.tasks[1:]
	}
}
