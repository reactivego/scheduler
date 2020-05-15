package scheduler

import (
	"fmt"
	"time"
)

// A task scheduled on an empty trampoline will dispatch sychronously
// and run immediately, while tasks scheduled by that task will dispatch
// asynchronously because they are added to a serial queue and executed at
// a later moment.
func Example_trampoline() {
	s := MakeTrampoline()

	fmt.Println("before")
	// Synchronous & Immediate
	s.Schedule(func() {
		fmt.Println("> outer")

		// Asynchronous & Serial
		s.Schedule(func() {
			fmt.Println("> inner")

			// Asynchronous & Serial
			s.Schedule(func() {
				fmt.Println("leaf")
			})

			fmt.Println("< inner")
		})

		fmt.Println("< outer")
	})
	s.Wait()
	fmt.Println("after")

	// Output:
	// before
	// > outer
	// < outer
	// > inner
	// < inner
	// leaf
	// after
}

func ExampleMakeTrampoline_scheduleRecursive() {
	tramp := MakeTrampoline()
	fmt.Println("before")

	i := 0
	tramp.ScheduleRecursive(func(self func()) {
		fmt.Println(i)
		i++
		if i < 3 {
			self()
		}
	})
	tramp.Wait()
	fmt.Println("after")

	// Output:
	// before
	// 0
	// 1
	// 2
	// after
}

func ExampleMakeTrampoline_scheduleFuture() {
	tramp := MakeTrampoline()
	fmt.Println("before")
	// Synchronous & Immediate
	tramp.ScheduleFuture(10*time.Millisecond, func() {
		fmt.Println("> outer")

		// Asynchronous & Serial
		tramp.Schedule(func() {
			fmt.Println("> inner")

			// Asynchronous & Serial
			tramp.Schedule(func() {
				fmt.Println("leaf")
			})

			fmt.Println("< inner")
		})

		fmt.Println("< outer")
	})
	tramp.Wait()
	fmt.Println("after")

	// Output:
	// before
	// > outer
	// < outer
	// > inner
	// < inner
	// leaf
	// after
}

func ExampleMakeTrampoline_scheduleFutureRecursive() {
	const asap = 0
	const _5ms = 5 * time.Millisecond
	const _10ms = 2 * _5ms
	const _20ms = 2 * _10ms

	tramp := MakeTrampoline()
	fmt.Println("before")

	tramp.ScheduleFutureRecursive(asap, func(self func(time.Duration)) {
		fmt.Println("> outer")

		//fmt.Println(time.Now().Sub(start).Round(_10ms))

		tramp.ScheduleFutureRecursive(_10ms, func(self func(time.Duration)) {
			fmt.Println("leaf 10ms")
		})

		tramp.ScheduleFutureRecursive(_5ms, func(self func(time.Duration)) {
			fmt.Println("leaf 5ms")
		})

		fmt.Println("< outer")
	})
	tramp.Wait()

	fmt.Println("after")

	// Output:
	// before
	// > outer
	// < outer
	// leaf 5ms
	// leaf 10ms
	// after
}

// Goroutine scheduler will dispatch a task asynchronously and run it
// concurrently with previously scheduled tasks. Nested tasks dispatched
// inside ScheduleRecursive by calling the function self() will be
// asynchronous and serial.
func Example_goroutine() {
	s := Goroutine

	fmt.Println("before")

	i := 0
	s.ScheduleRecursive(func(self func()) {
		fmt.Println(i)
		i++
		if i < 5 {
			self()
		}
	})
	fmt.Println("after")

	// Wait for the goroutine to finish.
	s.Wait()

	// Unordered output:
	// before
	// after
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleMakeGoroutine_cancel() {
	s := MakeGoroutine()

	const _10ms = 10 * time.Millisecond

	s.ScheduleFuture(_10ms, func() {
		// do nothing....
	})

	c := s.ScheduleFutureRecursive(_10ms, func(self func(due time.Duration)) {
		// do nothing....
		self(_10ms)
	})
	c.Cancel()

	time.Sleep(100 * time.Millisecond)

	fmt.Println(s)

	// Output:
	// Goroutine{ goroutines = 0 }
}
