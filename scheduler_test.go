package scheduler_test

import (
	"fmt"
	"time"

	"github.com/reactivego/scheduler"
)

// The concurrent Goroutine scheduler will dispatch a task asynchronously and
// run it concurrently with previously scheduled tasks. Nested tasks
// dispatched inside ScheduleRecursive by calling the function again() will be
// asynchronous and serial.
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
	// Output:
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
	serial := scheduler.NewTrampoline()

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

func ExampleNewTrampoline_scheduleRecursive() {
	serial := scheduler.NewTrampoline()

	i := 0
	serial.ScheduleRecursive(func(again func()) {
		fmt.Println(i)
		i++
		if i < 3 {
			again()
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

func ExampleNewTrampoline_scheduleLoop() {
	serial := scheduler.NewTrampoline()

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
	// Output:
	// BEFORE
	// 1
	// 2
	// 3
	// AFTER
	// tasks = 0
}

func ExampleNewTrampoline_scheduleFuture() {
	serial := scheduler.NewTrampoline()

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

func ExampleNewTrampoline_scheduleFutureRecursive() {
	const ms = time.Millisecond

	serial := scheduler.NewTrampoline()

	serial.ScheduleFutureRecursive(0*ms, func(again func(time.Duration)) {
		fmt.Println("> outer")

		serial.ScheduleFutureRecursive(10*ms, func(again func(time.Duration)) {
			fmt.Println("leaf 10ms")
		})

		serial.ScheduleFutureRecursive(5*ms, func(again func(time.Duration)) {
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

func ExampleGoroutine_cancel() {
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
	// Output:
	// tasks = 0
}
