# scheduler

    import "github.com/reactivego/scheduler"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/scheduler.svg)](https://pkg.go.dev/github.com/reactivego/scheduler#section-documentation)

Package `scheduler` provides a concurrent and a serial task scheduler with support for task cancellation.

The concurrent scheduler is exported as a global public variable with the name **`Goroutine`**.
This scheduler can be used directly. Alternatively, you can create a new concurrent scheduler
by calling **`NewConcurrentScheduler()`**.

A serial scheduler can be instantiated by calling either **`New()`** or **`NewSerialScheduler()`** function.

## Examples

### Concurrent

The concurrent Goroutine scheduler will dispatch tasks asynchronously and run them
concurrently with previously scheduled tasks. Nested tasks dispatched inside
ScheduleRecursive by calling the function again() will be added to a serial queue
and run in the order they were dispatched in.

Code:
```go
func Example_concurrent() {
	concurrent := scheduler.Goroutine

	i := 0
	concurrent.ScheduleRecursive(func(again func()) {
		fmt.Println(i)
		i++
		if i < 5 {
			again()
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

The serial scheduler will dispatch tasks asynchronously by adding
them to a serial queue and running them when the Wait method is called.

Code:
```go
func Example_serial() {
	serial := scheduler.New()

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

### Task Cancellation

You can cancel a scheduled task as shown in the example below:

Code:
```go
func Example_cancel() {
	const ms = time.Millisecond

	concurrent := scheduler.Goroutine

	concurrent.ScheduleFuture(10*ms, func() {
		// do nothing....
	})

	running := concurrent.ScheduleFutureRecursive(10*ms, func(again func(due time.Duration)) {
		// do nothing....
		again(10 * ms)
	})
	running.Cancel()

	concurrent.Wait()
	fmt.Println("tasks =", concurrent.Count())
}
```
Output:
```
tasks = 0
```

### Loop Scheduling

The ScheduleLoop method provides an easy way to implement loop-like behavior:

Code:
```go
func ExampleNew_scheduleLoop() {
	serial := scheduler.New()

	serial.ScheduleLoop(1, func(index int, again func(next int)) {
		fmt.Println(index)
		if index < 3 {
			again(index + 1)
		}
	})

	fmt.Println("BEFORE")
	serial.Wait()
	fmt.Println("AFTER")
	fmt.Println("tasks =", serial.Count())
}
```
Output:
```
BEFORE
1
2
3
AFTER
tasks = 0
```

## Interfaces

### Scheduler

Scheduler defines an interface for task execution management. Task scheduling happens
asynchronously without blocking the caller. Implementation may execute tasks
sequentially or concurrently.

```go
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

	// Wait will return when there are no more tasks running.
	Wait()

	// Gosched will give the scheduler an oportunity to run another task
	Gosched()

	// IsConcurrent returns true for a scheduler that runs tasks concurrently.
	// When using a concurrent scheduler, tasks will need to use synchronization
	// primitives like mutexes to properly guard against race conditions when
	// accessing shared data.
	IsConcurrent() bool

	// Count returns the number of currently active tasks.
	Count() int

	// String representation when printed.
	String() string
}
```

### SerialScheduler

SerialScheduler is a Scheduler that schedules tasks to run sequentially.
Tasks scheduled on this scheduler never access shared data at the same time.

```go
type SerialScheduler interface {
	Serial()
	Scheduler
}
```

### ConcurrentScheduler

ConcurrentScheduler is a Scheduler that schedules tasks concurrently.
Tasks will need to use synchronization primitives like mutexes to properly
guard against race conditions when accessing shared data.

```go
type ConcurrentScheduler interface {
	Concurrent()
	Scheduler
}
```

### Runner

Runner is an interface to a running task. It can be used to cancel the
running task by calling its Cancel() method.

```go
type Runner interface {
	// Cancel the running task.
	Cancel()
}
```
