package scheduler_test

import (
	"fmt"
	"time"

	"github.com/reactivego/scheduler"
)

// The concurrent Goroutine scheduler will dispatch a task asynchronously and
// run it concurrently with previously scheduled tasks. Nested tasks
// dispatched inside ScheduleRecursive by calling the function self() will be
// asynchronous and serial.
func Example_concurrent() {
	concurrent := scheduler.Goroutine

	fmt.Println("BEFORE")

	i := 0
	concurrent.ScheduleRecursive(func(self func()) {
		fmt.Println(i)
		i++
		if i < 5 {
			self()
		}
	})

	fmt.Println("AFTER")

	// Wait for the goroutine to finish.
	concurrent.Wait()
	fmt.Println("tasks =", concurrent.Count())
	// Unordered output:
	// BEFORE
	// AFTER
	// 0
	// 1
	// 2
	// 3
	// 4
	// tasks = 0
}

// The serial Trampoline scheduler will dispatch tasks asynchronously by adding
// them to a serial queue and running them when the Wait method is called.
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
	// Output:
	// BEFORE WAIT
	// > outer
	// < outer
	// > inner
	// < inner
	// leaf
	// AFTER WAIT (tasks = 0)
}

func ExampleMakeTrampoline_scheduleRecursive() {
	serial := scheduler.MakeTrampoline()

	i := 0
	serial.ScheduleRecursive(func(self func()) {
		fmt.Println(i)
		i++
		if i < 3 {
			self()
		}
	})

	fmt.Println("BEFORE")
	serial.Wait()
	fmt.Println("AFTER")
	fmt.Println("tasks =", serial.Count())
	// Output:
	// BEFORE
	// 0
	// 1
	// 2
	// AFTER
	// tasks = 0
}

func ExampleMakeTrampoline_scheduleFuture() {
	serial := scheduler.MakeTrampoline()

	// Asynchronous & Serial
	serial.ScheduleFuture(10*time.Millisecond, func() {
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
	// Output:
	// BEFORE WAIT
	// > outer
	// < outer
	// > inner
	// < inner
	// leaf
	// AFTER WAIT (tasks = 0)
}

func ExampleMakeTrampoline_scheduleFutureRecursive() {
	const ms = time.Millisecond

	serial := scheduler.MakeTrampoline()

	serial.ScheduleFutureRecursive(0*ms, func(self func(time.Duration)) {
		fmt.Println("> outer")

		serial.ScheduleFutureRecursive(10*ms, func(self func(time.Duration)) {
			fmt.Println("leaf 10ms")
		})

		serial.ScheduleFutureRecursive(5*ms, func(self func(time.Duration)) {
			fmt.Println("leaf 5ms")
		})

		fmt.Println("< outer")
	})

	fmt.Println("BEFORE WAIT")

	serial.Wait()

	fmt.Printf("AFTER WAIT (tasks = %d)\n", serial.Count())
	// Output:
	// BEFORE WAIT
	// > outer
	// < outer
	// leaf 5ms
	// leaf 10ms
	// AFTER WAIT (tasks = 0)
}

func ExampleMakeGoroutine_cancel() {
	const ms = time.Millisecond

	concurrent := scheduler.MakeGoroutine()

	concurrent.ScheduleFuture(10*ms, func() {
		// do nothing....
	})

	running := concurrent.ScheduleFutureRecursive(10*ms, func(self func(due time.Duration)) {
		// do nothing....
		self(10 * ms)
	})
	running.Cancel()

	concurrent.Wait()
	fmt.Println("tasks =", concurrent.Count())
	// Output:
	// tasks = 0
}
