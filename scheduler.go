package scheduler

import "time"

// Scheduler is an interface for scheduling tasks. Dispatching is either
// Synchronous or Asynchronous. Processing can be Immediate, Serial or
// Concurrent.
type Scheduler interface {
	// Now returns the current time according to the scheduler.
	Now() time.Time

	// Schedule dispatches a task to the scheduler.
	Schedule(task func())

	// ScheduleRecursive dispatches a task to the scheduler. Use the self
	// function to schedule another iteration of a repeating algorithm on
	// the scheduler.
	ScheduleRecursive(task func(self func()))

	// ScheduleFuture dispatches a task to the scheduler to be executed later.
	// The due time specifies when the task should be executed.
	ScheduleFuture(due time.Duration, task func())

	// ScheduleFutureRecursive dispatches a task to the scheduler to be
	// executed later. Use the self function to schedule another iteration of a
	// repeating algorithm on the scheduler. The due time specifies when the
	// task should be executed.
	ScheduleFutureRecursive(due time.Duration, task func(self func(time.Duration)))

	// Cancel will remove all queued tasks from the scheduler. A running task is
	// not affected by cancel and will continue until it is finished.
	Cancel()

	// IsAsynchronous returns true when the dispatch methods Schedule,
	// ScheduleRecursive, ScheduleFuture and ScheduleFutureRecursive
	// return before the scheduled task has run to completion.
	//
	// An empty Trampoline scheduler's first schedule method call will be
	// synchronous. This means the schedule methods will only return when the
	// task has finished running. For the Trampoline scheduler the value
	// returned here changes based on whether a task is currently scheduled on
	// the scheduler or not.
	IsAsynchronous() bool
}
