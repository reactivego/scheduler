package scheduler

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Goroutine is a concurrent scheduler. Schedule methods dispatch tasks 
// asynchronously, running them concurrently with previously scheduled tasks.
// It is safe to call the Goroutine scheduling methods from multiple
// concurrently running goroutines. Nested tasks dispatched inside e.g.
// ScheduleRecursive by calling the function self() will be added to a
// serial queue and run in the order they were dispatched in.
var Goroutine = MakeGoroutine()

// MakeGoroutine creates and returns a new concurrent scheduler instance.
// The returned instance implements the Scheduler interface.
func MakeGoroutine() *goroutine {
	return &goroutine{}
}

type goroutine struct {
	concurrent int32
}

func (s *goroutine) Now() time.Time {
	return time.Now()
}

// Schedule a task; dispatch it asynchronously to run concurrently as a
// new goroutine.
func (s *goroutine) Schedule(task func()) {
	go func() {
		atomic.AddInt32(&s.concurrent, 1)
		defer atomic.AddInt32(&s.concurrent, -1)
		task()
	}()
}

// ScheduleRecursive schedules a task; dispatches it asynchronously to run
// concurrently as a new goroutine. Inside a task scheduled with
// ScheduleRecursive, using the self() function will asynchronously
// reschedule the task to run concurrently with itself.
func (s *goroutine) ScheduleRecursive(task func(self func())) {
	go func() {
		atomic.AddInt32(&s.concurrent, 1)
		defer atomic.AddInt32(&s.concurrent, -1)
		MakeTrampoline().ScheduleRecursive(task)
	}()
}

// ScheduleFuture schedules a task; dispatches it asynchronously to run
// concurrently as a new goroutine. The goroutine will wait until the
// time is due to run the task.
func (s *goroutine) ScheduleFuture(due time.Duration, task func()) {
	go func() {
		atomic.AddInt32(&s.concurrent, 1)
		defer atomic.AddInt32(&s.concurrent, -1)
		time.Sleep(due)
		task()
	}()
}

// ScheduleFutureRecursive schedules a task; dispatches it asynchronously to
// run concurrently as a new goroutine at some moment in the future. Inside a
// task scheduled with ScheduleRecursiveFuture, using the self(due) function
// will asynchronously reschedule the task to run concurrently with itself
// at some moment in time in the future.
func (s *goroutine) ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration))) {
	go func() {
		atomic.AddInt32(&s.concurrent, 1)
		defer atomic.AddInt32(&s.concurrent, -1)
		MakeTrampoline().ScheduleFutureRecursive(due, task)
	}()
}


// Cancel will remove all queued tasks from the scheduler. A running task is
// not affected by cancel and will continue until it is finished.
func (s *goroutine) Cancel() {
}

// IsAsynchronous returns true.
func (s *goroutine) IsAsynchronous() bool {
	return true
}

func (s *goroutine) String() string {
	return fmt.Sprintf("Goroutine{ Asynchronous:Concurrent(%d) }", s.concurrent)
}
