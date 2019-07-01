package scheduler

import "time"

// Scheduler is an interface for scheduling tasks. Dispatching is either
// Synchronous or Asynchronous. Processing can be Immediate, Serial or
// Concurrent.
type Scheduler interface {
	// Schedule dispatches a task to the scheduler.
	Schedule(task func())

	// ScheduleRecursive dispatches a task to the scheduler. Use the self
	// function to schedule another iteration of a repeating algorithm on
	// the scheduler.
	ScheduleRecursive(task func(self func()))

	// ScheduleFutureRecursive will dispatch the first task synchronously
	// and any subsequent tasks asynchronously on a task queue. When the first
	// task eventually returns the queue of tasks is empty again.
	// The task will be scheduled after the time period has passed.
	// To schedule for the next period, the code should call the self function.
	// Just returning from the task function will terminate a task.
	ScheduleFutureRecursive(timeout time.Duration, task func(self func(time.Duration)))

	// IsAsynchronous returns true when the dispatch methods Schedule
	// and ScheduleRecursive methods return before the scheduled task
	// has run to completion. For some schedulers the value returned
	// here changes based on whether a task is currently scheduled on
	// the scheduler.
	IsAsynchronous() bool

	// IsSerial returns true when the scheduler adds a scheduled task to
	// a queue. A single goroutine then takes tasks of the queue and runs
	// them in sequence. For some schedulers the value returned here
	// changes based on whether a task is currently scheduled on
	// the scheduler. A scheduler that is not serial or concurrent
	// runs tasks immediately.
	IsSerial() bool

	// IsConcurrent returns true when the scheduler starts a scheduled
	// task to run concurrently alongside other scheduled tasks. A scheduler
	// that is not serial or concurrent runs tasks immediately.
	IsConcurrent() bool
}
