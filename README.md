# scheduler

    import "github.com/reactivego/scheduler"

[![](svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/scheduler?tab=doc)
[![](svg/godoc.svg)](https://godoc.org/github.com/reactivego/scheduler)

Package `scheduler` provides a concurrent and a serial task scheduler with support for task cancellation.

The concurrent scheduler is exported as a global public variable with the name **`Goroutine`**.
This scheduler can be used directly.

A serial scheduler needs to be instantiated by calling the **`MakeTrampoline`** function exported by this package.

## Examples

### Concurrent

The concurrent Goroutine scheduler will dispatch a task by running it
concurrently with previously scheduled tasks. These may start running
immediately after they have been scheduled. Nested tasks dispatched by calling
the self() function will be placed on a task queue and run in sequence after
the currently scheduled task returns.

Code:
```go
func Example_concurrent() {
	concurrent := scheduler.Goroutine

	i := 0
	concurrent.ScheduleRecursive(func(self func()) {
		fmt.Println(i)
		i++
		if i < 5 {
			self()
		}
	})

	// Wait for the goroutine to finish.
	concurrent.Wait()
	fmt.Println("tasks =", concurrent.Count())
}
```
Output:
```
0
1
2
3
4
tasks = 0
```

### Serial

The serial Trampoline scheduler will dispatch tasks by adding them to a serial
queue and running them when the Wait method is called on the scheduler.

Code:
```go
func Example_serial() {
	serial := scheduler.MakeTrampoline()

	// Asynchronous & serial
	serial.Schedule(func() {
		fmt.Println("> outer")

		// Asynchronous & Serial
		serial.Schedule(func() {
			fmt.Println("> inner")

			// Asynchronous & Serial
			serial.Schedule(func() {
				fmt.Println("leaf")
			})

			fmt.Println("< inner")
		})

		fmt.Println("< outer")
	})

	fmt.Println("BEFORE WAIT")

	serial.Wait()

	fmt.Printf("AFTER WAIT (tasks = %d)\n", serial.Count())
}
```
Output:
```
BEFORE WAIT
> outer
< outer
> inner
< inner
leaf
AFTER WAIT (tasks = 0)
```

## Interfaces

### Scheduler 

Scheduler is an interface for running tasks. Scheduling of tasks is
asynchronous/non-blocking. Tasks can be executed in sequence or concurrently.

```go
type Scheduler interface {
	// Now returns the current time according to the scheduler.
	Now() time.Time

	// Since returns the time elapsed, is a shorthand for Now().Sub(t).
	Since(t time.Time) time.Duration

	// Schedule dispatches a task to the scheduler.
	Schedule(task func()) Runner

	// ScheduleRecursive dispatches a task to the scheduler. Use the self
	// function to schedule another iteration of a repeating algorithm on
	// the scheduler.
	ScheduleRecursive(task func(self func())) Runner

	// ScheduleFuture dispatches a task to the scheduler to be executed later.
	// The due time specifies when the task should be executed.
	ScheduleFuture(due time.Duration, task func()) Runner

	// ScheduleFutureRecursive dispatches a task to the scheduler to be
	// executed later. Use the self function to schedule another iteration of a
	// repeating algorithm on the scheduler. The due time specifies when the
	// task should be executed.
	ScheduleFutureRecursive(due time.Duration, task func(self func(due time.Duration))) Runner

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
```
### Runner

Runner is an interface to a running task. It can be used to cancel the running
task by calling its Cancel() method.

```go
type Runner interface {
	// Cancel the running task.
	Cancel()
}
```
