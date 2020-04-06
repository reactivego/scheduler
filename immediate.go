package scheduler

import "time"

// Immediate scheduler will dispatch a task synchronously and run it
// immediately. It will also schedule recursive tasks immediately,
// so it can run out of stack space for very deep recursion.
// It is safe to use the Immediate scheduler from multiple concurrently
// running goroutines.
var Immediate = &immediate{}

type immediate struct{}

func (s immediate) Now() time.Time {
	return time.Now()
}

func (s immediate) Schedule(task func()) {
	task()
}

func (s immediate) ScheduleRecursive(task func(self func())) {
	task(func() { s.ScheduleRecursive(task) })
}

func (s immediate) ScheduleFuture(due time.Duration, task func()) {
	time.Sleep(due)
	task()
}

func (s immediate) ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration))) {
	time.Sleep(due)
	task(func(due time.Duration) { s.ScheduleFutureRecursive(due, task) })
}

func (s immediate) Cancel() {
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

func (s immediate) String() string {
	return "Immediate{ Synchronous:Immediate }"
}
