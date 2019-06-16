# scheduler

    import "github.com/reactivego/scheduler"

[![](https://godoc.org/github.com/reactivego/scheduler?status.png)](http://godoc.org/github.com/reactivego/scheduler)

Package `scheduler` contains different implementations of the `Scheduler` interface. Scheduling can be characterized in two ways. First, in the way tasks are dispatched to the scheduler, which can be either *asynchronous* or *synchronous*. Second, in the way tasks are actually executed by the *scheduler*, which can be *immediate*, *serial* or *concurrent*.

**Root Dispatch** is defined as using the `Schedule` or `ScheduleRecursive` method of a scheduler to schedule a task.

**Nested Dispatch** is defined as using the `self` function inside `ScheduleRecursive` to schedule a nested task.

Below is an overview of the schedulers exported by the `scheduler` package:


| Scheduler 			| Root Dispatch			| Nested Dispatch		|
| ---:      				| ---  						| ---    					|
| **`Immediate`**			| synchronous & immediate 	| synchronous & immediate 	|
| **`CurrentGoroutine`**	| synchronous & immediate <sup>*</sup> 	| asynchronous & serial		|
| *`Trampoline`*				| synchronous & immediate <sup>*</sup> 	| asynchronous & serial		|
| **`NewGoroutine`**		| asynchronous & concurrent	| asynchronous & serial		|
| *`Goroutine`*				| asynchronous & concurrent	| asynchronous & concurrent	|

> <sup>*</sup> The scheduler is *synchronous & immediate* when no tasks are scheduled; however, when at least one task is scheduled it is *asynchronous & serial*.

The schedulers **`Immediate`**, **`CurrentGoroutine`** and **`NewGoroutine`** are predefined variables and can be used directly. The other two schedulers *`Trampoline`* and *`Goroutine`* both of type *`struct`* need to be instantiated, e.g. `scheduler := &Trampoline{}`, before they can be used.
