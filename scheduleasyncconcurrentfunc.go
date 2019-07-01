package scheduler

import "time"

// ScheduleAsyncConcurrentFunc is a function that can dispatch tasks
// asynchronously and then run them concurrently. The root scheduler
// is asynchronous and concurrent, while the ScheduleRecursive method
// creates a Trampoline to make recursive scheduling asynchronous and
// serial.
type ScheduleAsyncConcurrentFunc func(task func())

// Schedule the task, dispatching it asynchronously and running it
// concurrently with other scheduled tasks.
func (s ScheduleAsyncConcurrentFunc) Schedule(task func()) {
	s(task)
}

// ScheduleRecursive schedules the task asynchronous and
// concurrent. It creates a dedicated Trampoline scheduler
// for the task so calling self inside the task will schedule a
// task asynchronous and serial.
func (s ScheduleAsyncConcurrentFunc) ScheduleRecursive(task func(self func())) {
	inner := &Trampoline{}
	s(func() { inner.ScheduleRecursive(task) })
}

func (s ScheduleAsyncConcurrentFunc) ScheduleFutureRecursive(timeout time.Duration, task func(self func(time.Duration))) {
	inner := &Trampoline{}
	s(func() { inner.ScheduleFutureRecursive(timeout,task) })
}

// IsAsynchronous returns true.
func (s ScheduleAsyncConcurrentFunc) IsAsynchronous() bool {
	return true
}

// IsSerial returns false.
func (s ScheduleAsyncConcurrentFunc) IsSerial() bool {
	return false
}

// IsConcurrent returns true.
func (s ScheduleAsyncConcurrentFunc) IsConcurrent() bool {
	return true
}

// Wait does nothing for this scheduler.
func (s ScheduleAsyncConcurrentFunc) Wait(onCancel func(func())) {
}
