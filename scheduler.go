// Package scheduler provides a concurrent and a serial task scheduler with
// support for task cancellation.
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

// SerialScheduler is a Scheduler that schedules tasks to run sequentially.
// Tasks scheduled on this scheduler never access shared data at the same time.
type SerialScheduler interface {
	Serial()
	Scheduler
}

// New creates and returns a serial (non-concurrent) scheduler that runs all
// tasks on a single goroutine. The returned scheduler is returned as a
// SerialScheduler interface. Tasks scheduled will be dispatched asynchronously
// because they are added to a serial queue. When the Wait method is called all
// tasks scheduled on the serial queue will be performed in the same order in
// which they were added to the queue.
//
// The returned scheduler is not safe to be shared by multiple goroutines
// concurrently. It should be used purely from a single goroutine to schedule
// tasks to run sequentially.
func New() SerialScheduler {
	return &trampoline{}
}

func NewSerialScheduler() SerialScheduler {
	return &trampoline{}
}

// ConcurrentScheduler is a Scheduler that schedules tasks concurrently.
// Tasks scheduled on this scheduler may access shared data at the same time.
type ConcurrentScheduler interface {
	Concurrent()
	Scheduler
}

// Goroutine is a concurrent scheduler. Schedule methods dispatch tasks
// asynchronously, running them concurrently with previously scheduled tasks.
// It is safe to call the Goroutine scheduling methods from multiple
// concurrently running goroutines. Nested tasks dispatched inside e.g.
// ScheduleRecursive by calling the function again() will be added to a
// serial queue and run in the order they were dispatched in.
var Goroutine = ConcurrentScheduler(&goroutine{})

func NewConcurrentScheduler() ConcurrentScheduler {
	return &goroutine{}
}
