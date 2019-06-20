// Package scheduler implements task schedulers.
//
// Tasks can be dispatched asynchronously or synchronously to a scheduler.
// Asynchronous means the dispatch function returns before the task starts,
// whereas synchronous means the dispatch function returns only afer the
// dispatched task has completed.
// 
// A scheduler runs dispatched tasks concurrently, sequentially (serial)
// or immediately.
package scheduler
