package scheduler

import "time"

// ScheduleAsyncSerialFunc is a function that can dispatch tasks asynchronously
// and run them in sequence. The root scheduler as well as recursive scheduling
// is asynchronous and serial.
//
// An async/serial scheduler is usually paired with a sync/serial scheduler
// using the same serial task queue.
// The async/serial task dispatched from a goroutine deposits the result of
// some work done in the background.
// The sync/serial task dispatched from the main goroutine copies this result
// out to a local variable to be processed further.
type ScheduleAsyncSerialFunc func(task func())

// Schedule the task, dispatching it asynchronously on a serial queue.
func (s ScheduleAsyncSerialFunc) Schedule(task func()) {
	s(task)
}

// Schedule the task and recursive tasks, dispatching it asynchronously on
// a serial queue.
func (s ScheduleAsyncSerialFunc) ScheduleRecursive(task func(self func())) {
	self := func() {
		s.ScheduleRecursive(task)
	}
	s(func() {
		task(self)
	})
}

func (s ScheduleAsyncSerialFunc) ScheduleFutureRecursive(timeout time.Duration, task func(self func(time.Duration))) {
	self := func(timeout time.Duration) {
		s.ScheduleFutureRecursive(timeout, task)
	}
	s(func() {
		time.Sleep(timeout)
		task(self)
	})
}

// IsAsynchronous returns true.
func (s ScheduleAsyncSerialFunc) IsAsynchronous() bool {
	return true
}

// IsSerial returns true.
func (s ScheduleAsyncSerialFunc) IsSerial() bool {
	return true
}

// IsConcurrent returns false.
func (s ScheduleAsyncSerialFunc) IsConcurrent() bool {
	return false
}

// Wait does nothing for this scheduler.
func (s ScheduleAsyncSerialFunc) Wait(onCancel func(func())) {
}
