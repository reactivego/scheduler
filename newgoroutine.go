package scheduler

import (
	"time"
)

// NewGoroutine scheduler will dispatch a task asynchronously and run it
// concurrently with previously scheduled tasks. It is safe to call the
// NewGoroutine scheduler from multiple concurrently running goroutines.
// Nested tasks dispatched inside ScheduleRecursive by calling the
// function self() will be asynchronous and serial.
var NewGoroutine = &newgoroutine{}

type newgoroutine struct{}

func (s newgoroutine) Now() time.Time {
	return time.Now()
}

func (s newgoroutine) Schedule(task func()) {
	go task()
}

func (s newgoroutine) ScheduleRecursive(task func(self func())) {
	inner := &Trampoline{}
	go inner.ScheduleRecursive(task)
}

func (s newgoroutine) ScheduleFuture(due time.Duration, task func()) {
	go func() {
		time.Sleep(due)
		task()
	}()
}

func (s newgoroutine) ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration))) {
	inner := &Trampoline{}
	go inner.ScheduleFutureRecursive(due, task)
}

func (s newgoroutine) IsAsynchronous() bool {
	return true
}

func (s newgoroutine) IsSerial() bool {
	return false
}

func (s newgoroutine) IsConcurrent() bool {
	return true
}
