package scheduler

import (
	"fmt"
	"sync"
	"time"
)

// Immediate scheduler will dispatch a task synchronously and run it
// immediately. It will also schedule recursive tasks immediately,
// so it can run out of stack space for very deep recursion.
func Example_immediate() {
	// Synchronous & Immediate
	fmt.Println("before")
	Immediate.Schedule(func() {
		fmt.Println("> outer")

		Immediate.Schedule(func() {
			fmt.Println("> inner")

			Immediate.Schedule(func() {
				fmt.Println("leaf")
			})

			fmt.Println("< inner")
		})

		fmt.Println("< outer")
	})
	fmt.Println("after")

	// Output:
	// before
	// > outer
	// > inner
	// leaf
	// < inner
	// < outer
	// after
}

// CurrentGoroutine scheduler is a Trampoline scheduler. A task scheduled on
// an empty trampoline will be dispatched sychronously and run immediately,
// while tasks scheduled by that task will be asynchronous and serial.
func Example_currentGoroutine() {
	fmt.Println("before")
	// Synchronous & Immediate
	CurrentGoroutine.Schedule(func() {
		fmt.Println("> outer")

		// Asynchronous & Serial
		CurrentGoroutine.Schedule(func() {
			fmt.Println("> inner")

			// Asynchronous & Serial
			CurrentGoroutine.Schedule(func() {
				fmt.Println("leaf")
			})

			fmt.Println("< inner")
		})

		fmt.Println("< outer")
	})
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

// NewGoroutine scheduler will dispatch a task asynchronously and run it
// concurrently with previously scheduled tasks. Nested tasks dispatched
// inside ScheduleRecursive by calling the function self() will be
// asynchronous and serial.
func Example_newGoroutine() {
	fmt.Println("before")

	var wg sync.WaitGroup
	wg.Add(1)
	i := 0
	NewGoroutine.ScheduleRecursive(func(self func()) {
		fmt.Println(i)
		i++
		if i < 5 {
			self()
		} else {
			wg.Done()
		}
	})
	fmt.Println("after")

	// Wait for the goroutine to finish.
	wg.Wait()

	// Output:
	// before
	// after
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleTrampoline() {
	tramp := MakeTrampoline()
	fmt.Println("before")
	// Synchronous & Immediate
	tramp.Schedule(func() {
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

func ExampleTrampoline_ScheduleRecursive() {
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
	fmt.Println("after")

	// Output:
	// before
	// 0
	// 1
	// 2
	// after
}

func ExampleTrampoline_ScheduleFuture() {
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

func ExampleTrampoline_ScheduleFutureRecursive() {
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

	fmt.Println("after")

	// Output:
	// before
	// > outer
	// < outer
	// leaf 5ms
	// leaf 10ms
	// after
}

func ExampleNewGoroutine_cancel() {
	s := NewGoroutine

	const _10ms = 10 * time.Millisecond

	s.ScheduleFuture(_10ms, func() {
		// do nothing....
	})

	s.ScheduleFutureRecursive(_10ms, func(self func(due time.Duration)) {
		// do nothing....
		self(_10ms)
	})

	time.Sleep(100 * time.Millisecond)

	s.Cancel()
	fmt.Println(s)

	// Output:
	// NewGoroutine{ Asynchronous:Concurrent(0) }
}
