package scheduler

// ScheduleFunc is a function that can schedule tasks.
// The root scheduler as well as recursive scheduling is synchronous and immediate.
type ScheduleFunc func(task func())

// Schedule the task to run synchronously and immediate.
func (s ScheduleFunc) Schedule(task func()) {
	s(task)
}

// Schedule the task and recursive tasks to run synchronously and immediate.
func (s ScheduleFunc) ScheduleRecursive(task func(self func())) {
	self := func() { s.ScheduleRecursive(task) }
	s(func() { task(self) })
}

// IsAsynchronous returns false.
func (s ScheduleFunc) IsAsynchronous() bool {
	return false
}

// IsSerial returns false.
func (s ScheduleFunc) IsSerial() bool {
	return false
}

// IsConcurrent returns false.
func (s ScheduleFunc) IsConcurrent() bool {
	return false
}

// Wait does nothing for the ScheduleFunc scheduler.
func (s ScheduleFunc) Wait(onCancel func(func())) {
}
