package scheduler

import (
	"fmt"
	"sync/atomic"
	"time"
)

// NewGoroutine scheduler will dispatch a task asynchronously and run it
// concurrently with previously scheduled tasks. It is safe to call the
// NewGoroutine scheduler from multiple concurrently running goroutines.
// Nested tasks dispatched inside ScheduleRecursive by calling the
// function self() will be asynchronous and serial.
var NewGoroutine = makeNewGoroutine()

func makeNewGoroutine() *newgoroutine {
	return &newgoroutine{}
}

type newgoroutine struct {
	concurrent int32
}

func (s *newgoroutine) Now() time.Time {
	return time.Now()
}

func (s *newgoroutine) Schedule(task func()) {
	go func() {
		atomic.AddInt32(&s.concurrent, 1)
		defer atomic.AddInt32(&s.concurrent, -1)
		task()
	}()
}

func (s *newgoroutine) ScheduleRecursive(task func(self func())) {
	go func() {
		atomic.AddInt32(&s.concurrent, 1)
		defer atomic.AddInt32(&s.concurrent, -1)
		MakeTrampoline().ScheduleRecursive(task)
	}()
}

func (s *newgoroutine) ScheduleFuture(due time.Duration, task func()) {
	go func() {
		atomic.AddInt32(&s.concurrent, 1)
		defer atomic.AddInt32(&s.concurrent, -1)
		time.Sleep(due)
		task()
	}()
}

func (s *newgoroutine) ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration))) {
	go func() {
		atomic.AddInt32(&s.concurrent, 1)
		defer atomic.AddInt32(&s.concurrent, -1)
		MakeTrampoline().ScheduleFutureRecursive(due, task)
	}()
}

func (s *newgoroutine) Cancel() {
}

func (s *newgoroutine) IsAsynchronous() bool {
	return true
}

func (s *newgoroutine) IsSerial() bool {
	return false
}

func (s *newgoroutine) IsConcurrent() bool {
	return true
}

func (s *newgoroutine) String() string {
	return fmt.Sprintf("NewGoroutine{ Asynchronous:Concurrent(%d) }", s.concurrent)
}
