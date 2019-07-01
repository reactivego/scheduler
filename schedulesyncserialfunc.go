package scheduler

import "time"

// ScheduleSyncSerialFunc is a function that can dispatch tasks synchronously
// and run them in sequence. The root scheduler as well as recursive scheduling
// is synchronous and serial.
//
// A sync/serial scheduler is usually paired with an async/serial scheduler
// using the same serial task queue.
// The async/serial task dispatched from a goroutine deposits the result of
// some work done in the background.
// The sync/serial task dispatched from the main goroutine copies this result
// out to a local variable to be processed further.
type ScheduleSyncSerialFunc func(task func())

// Schedule the task, dispatching it synchronously on a serial queue.
func (s ScheduleSyncSerialFunc) Schedule(task func()) {
	s(task)
}

// Schedule the task and recursive tasks, dispatching it synchronously on
// a serial queue.
func (s ScheduleSyncSerialFunc) ScheduleRecursive(task func(self func())) {
	self := func() {
		s.ScheduleRecursive(task)
	}
	s(func() {
		task(self)
	})
}

func (s ScheduleSyncSerialFunc) ScheduleFutureRecursive(timeout time.Duration, task func(self func(time.Duration))) {
	self := func(timeout time.Duration) {
		s.ScheduleFutureRecursive(timeout, task)
	}
	s(func() {
		time.Sleep(timeout)
		task(self)
	})
}

// IsAsynchronous returns false.
func (s ScheduleSyncSerialFunc) IsAsynchronous() bool {
	return false
}

// IsSerial returns true.
func (s ScheduleSyncSerialFunc) IsSerial() bool {
	return true
}

// IsConcurrent returns false.
func (s ScheduleSyncSerialFunc) IsConcurrent() bool {
	return false
}

// Wait does nothing for this scheduler.
func (s ScheduleSyncSerialFunc) Wait(onCancel func(func())) {
}
