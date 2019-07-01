package scheduler

import (
	"sync"
	"time"
)

// NewGoroutine scheduler will dispatch a task asynchronously and run it
// concurrently with previously scheduled tasks. It is safe to call the
// NewGoroutine scheduler from multiple concurrently running goroutines.
// Nested tasks dispatched inside ScheduleRecursive by calling the
// function self() will be asynchronous and serial.
var NewGoroutine = &newgoroutine{}

type newgoroutine struct{}

func (s newgoroutine) Schedule(task func()) {
	go task()
}

func (s newgoroutine) ScheduleRecursive(task func(self func())) {
	inner := &Trampoline{}
	go inner.ScheduleRecursive(task)
}

func (s newgoroutine) ScheduleFutureRecursive(timeout time.Duration, task func(self func(time.Duration))) {
	inner := &Trampoline{}
	go inner.ScheduleFutureRecursive(timeout,task)
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

func (s newgoroutine) Wait(onCancel func(func())) {
	var wg sync.WaitGroup
	wg.Add(1)
	onCancel(wg.Done)
	wg.Wait()
}
