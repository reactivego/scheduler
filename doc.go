// Package scheduler provides task schedulers.
//
// Schedule a task to be executed using one of the dispatch methods:
// Schedule(task), ScheduleRecursive(task), ScheduleFuture(due, task)
// or ScheduleFutureRecursive(due,task)
//
// Task dispatch can be either asynchronous or synchronous. Asynchronous
// dispatch means the schedule method returns before the task
// starts, whereas synchronous dispatch means the schedule method returns
// only afer the dispatched task has completed.
package scheduler
