package scheduler

import (
	"time"
)

// Goroutine is a scheduler that dispatches tasks asynchronously and runs them
// concurrently. It is safe to use the Goroutine scheduler from concurrently
// running tasks.
// Note that the recursive scheduling functions will schedule recursively on a
// new goroutine.
type Goroutine struct{}

func MakeGoroutine() *Goroutine {
	return &Goroutine{}
}

func (s Goroutine) Now() time.Time {
	return time.Now()
}

// Schedule a task; dispatch it asynchronously to run concurrently as a
// new goroutine.
func (s Goroutine) Schedule(task func()) {
	go task()
}

// ScheduleRecursive schedules a task; dispatches it asynchronously to run
// concurrently as a new goroutine. Inside a task scheduled with
// ScheduleRecursive, using the self() function will asynchronously
// reschedule the task to run concurrently with itself.
func (s Goroutine) ScheduleRecursive(task func(self func())) {
	go task(func() { s.ScheduleRecursive(task) })
}

// ScheduleFuture schedules a task; dispatches it asynchronously to run
// concurrently as a new goroutine. The goroutine will wait until the
// time is due to run the task.
func (s Goroutine) ScheduleFuture(due time.Duration, task func()) {
	go func() {
		time.Sleep(due)
		task()
	}()
}

// ScheduleFutureRecursive schedules a task; dispatches it asynchronously to
// run concurrently as a new goroutine at some moment in the future. Inside a
// task scheduled with ScheduleRecursiveFuture, using the self(due) function
// will asynchronously reschedule the task to run concurrently with itself
// at some moment in time in the future.
func (s Goroutine) ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration))) {
	self := func(timeout time.Duration) {
		s.ScheduleFutureRecursive(due, task)
	}
	go func() {
		time.Sleep(due)
		task(self)
	}()
}

// Cancel will remove all queued tasks from the scheduler. A running task is
// not affected by cancel and will continue until it is finished.
func (s Goroutine) Cancel() {
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

func (s Goroutine) String() string {
	return "Goroutine{ Asynchronous:Concurrent }"
}
