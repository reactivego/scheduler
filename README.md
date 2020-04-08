# scheduler

    import "github.com/reactivego/scheduler"

[![](https://godoc.org/github.com/reactivego/scheduler?status.png)](http://godoc.org/github.com/reactivego/scheduler)

Package `scheduler` implements task schedulers. Scheduling can be characterized in two ways. First, in the way tasks are dispatched to the scheduler, which can be either *asynchronous* or *synchronous*. Second, in the way tasks are actually executed by the *scheduler*, which can be *immediate*, *serial* or *concurrent*.

**Root Dispatch** is defined as using the `Schedule` or `ScheduleRecursive` method of a scheduler to schedule a task.

**Recursive Dispatch** is defined as using the `self` function inside `ScheduleRecursive` to schedule a nested task.

This package defines the schedulers **`Trampoline`** and **`Goroutine`** as public variables that can be used directly.

To manually create a scheduler, use e.g. `scheduler.MakeTrampoline()` or `scheduler.MakeGoroutine()`.
