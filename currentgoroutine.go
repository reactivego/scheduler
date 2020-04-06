package scheduler

// CurrentGoroutine scheduler is a Trampoline scheduler. A task scheduled on
// an empty trampoline will be dispatched sychronously and run immediately,
// while tasks scheduled by that task will be asynchronous and serial.
// The CurrentGoroutine scheduler is not safe to use from multiple goroutines
// at the same time. It should be used purely for scheduling tasks from the
// current goroutine.
var CurrentGoroutine = MakeTrampoline()
