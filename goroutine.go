package scheduler

import (
	"sync"
	"time"
)

// Goroutine is a scheduler that dispatches tasks asynchronously and runs them
// concurrently. It is safe to use the Goroutine scheduler from concurrently
// running tasks.
type Goroutine struct{}

// Schedule a task; dispatch it asynchronously to run concurrently as a
// new goroutine.
func (s Goroutine) Schedule(task func()) {
	go task()
}

// Schedule a task; dispatch it asynchronously to run concurrently as a
// new goroutine. Inside a task scheduled with ScheduleRecursive, using the
// self() function will asynchronously reschedule the task to run concurrently
// with itself.
func (s Goroutine) ScheduleRecursive(task func(self func())) {
	go task(func() { s.ScheduleRecursive(task) })
}

func (s Goroutine) ScheduleFutureRecursive(timeout time.Duration, task func(self func(time.Duration))) {
	self := func(timeout time.Duration) {
		s.ScheduleFutureRecursive(timeout, task)
	}
	go func() {
		time.Sleep(timeout)
		task(self)
	}()
}

// IsAsynchronous returns true.
func (s Goroutine) IsAsynchronous() bool {
	return true
}

// IsSerial returns false.
func (s Goroutine) IsSerial() bool {
	return false
}

// IsConcurrent returns true.
func (s Goroutine) IsConcurrent() bool {
	return true
}

