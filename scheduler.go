package scheduler

import "time"

// Scheduler is an interface for running tasks.
// Scheduling of tasks is asynchronous/non-blocking.
// Tasks can be executed in sequence or concurrently.
type Scheduler interface {
	// Now returns the current time according to the scheduler.
	Now() time.Time

	// Since returns the time elapsed, is a shorthand for Now().Sub(t).
	Since(t time.Time) time.Duration

	// Schedule dispatches a task to the scheduler.
	Schedule(task func()) Runner

	// ScheduleRecursive dispatches a task to the scheduler. Use the again
	// function to schedule another iteration of a repeating algorithm on
	// the scheduler.
	ScheduleRecursive(task func(again func())) Runner

	// ScheduleLoop dispatches a task to the scheduler. Use the again
	// function to schedule another iteration of a repeating algorithm on
	// the scheduler. The current loop index is passed to the task. The loop
	// index starts at the value passed in the from argument. The task is
	// expected to pass the next loop index to the again function.
	ScheduleLoop(from int, task func(index int, again func(next int))) Runner

	// ScheduleFuture dispatches a task to the scheduler to be executed later.
	// The due time specifies when the task should be executed.
	ScheduleFuture(due time.Duration, task func()) Runner

	// ScheduleFutureRecursive dispatches a task to the scheduler to be
	// executed later. Use the again function to schedule another iteration of a
	// repeating algorithm on the scheduler. The due time specifies when the
	// task should be executed.
	ScheduleFutureRecursive(due time.Duration, task func(again func(due time.Duration))) Runner

	// Wait will return when the Cancel() method is called or when there are no
	// more tasks running. Note, the currently running task may schedule
	// additional tasks to the queue to run later.
	Wait()

	// Gosched will give the scheduler an oportunity to run another task
	Gosched()

	// IsConcurrent returns true for a scheduler that runs tasks concurrently.
	IsConcurrent() bool

	// Count returns the number of currently active tasks.
	Count() int

	// String representation when printed.
	String() string
}

// Runner is an interface to a running task. It can be used to cancel the
// running task by calling its Cancel() method.
type Runner interface {
	// Cancel the running task.
	Cancel()
}
