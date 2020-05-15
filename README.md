# scheduler

    import "github.com/reactivego/scheduler"

[![](https://godoc.org/github.com/reactivego/scheduler?status.png)](http://godoc.org/github.com/reactivego/scheduler)

Package `scheduler` implements task schedulers. Tasks are scheduled in a non-blocking asynchronous way.
Depending on the scheduler the tasks are either executed *in-sequence* or *concurrently*.

**Root Dispatch** is defined as using the `Schedule` or `ScheduleRecursive` or `ScheduleFutureRecursive` method of a scheduler to schedule a task.

**Recursive Dispatch** is defined as using the `self` function inside `ScheduleRecursive` or `ScheduleFutureRecursive` to schedule a nested task.

This package defines the scheduler **`Goroutine`** as a public variable that can be used directly.

To manually create a scheduler, use e.g. `scheduler.MakeTrampoline()` or `scheduler.MakeGoroutine()`.
