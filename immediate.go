package scheduler

import "time"

// Immediate scheduler will dispatch a task synchronously and run it
// immediately. It will also schedule recursive tasks immediately,
// so it can run out of stack space for very deep recursion.
// It is safe to use the Immediate scheduler from multiple concurrently
// running goroutines.
var Immediate = &immediate{}

type immediate struct{}

func (s immediate) Schedule(task func()) {
	task()
}

func (s immediate) ScheduleRecursive(task func(self func())) {
	self := func() { s.ScheduleRecursive(task) }
	task(self)
}

func (s immediate) ScheduleFutureRecursive(timeout time.Duration, task func(self func(time.Duration))) {
	self := func(timeout time.Duration) {
		s.ScheduleFutureRecursive(timeout, task)
	}
	time.Sleep(timeout)
	task(self)
}

func (s immediate) IsAsynchronous() bool {
	return false
}

func (s immediate) IsSerial() bool {
	return false
}

func (s immediate) IsConcurrent() bool {
	return false
}

func (s immediate) Wait(onCancel func(func())) {
}
