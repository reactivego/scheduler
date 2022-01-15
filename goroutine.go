package scheduler

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Goroutine is a concurrent scheduler. Schedule methods dispatch tasks
// asynchronously, running them concurrently with previously scheduled tasks.
// It is safe to call the Goroutine scheduling methods from multiple
// concurrently running goroutines. Nested tasks dispatched inside e.g.
// ScheduleRecursive by calling the function again() will be added to a
// serial queue and run in the order they were dispatched in.
var Goroutine = &goroutine{}

// cancel

type cancel chan struct{}

func (c cancel) Cancel() {
	close(c)
}

// goroutine

type goroutine struct {
	sync.Mutex
	concurrent sync.WaitGroup
	active     int32
}

func (s *goroutine) Now() time.Time {
	return time.Now()
}

func (s *goroutine) Since(t time.Time) time.Duration {
	return s.Now().Sub(t)
}

func (s *goroutine) Schedule(task func()) Runner {
	cancel := make(cancel)
	atomic.AddInt32(&s.active, 1)
	s.concurrent.Add(1)
	go func() {
		defer atomic.AddInt32(&s.active, -1)
		defer s.concurrent.Done()
		select {
		case <-cancel:
			// cancel
		default:
			task()
		}
	}()
	return cancel
}

func (s *goroutine) ScheduleRecursive(task func(again func())) Runner {
	runner := make(chan Runner, 1)
	atomic.AddInt32(&s.active, 1)
	s.concurrent.Add(1)
	go func() {
		defer atomic.AddInt32(&s.active, -1)
		defer s.concurrent.Done()
		serial := New()
		runner <- serial.ScheduleRecursive(task)
		serial.Wait()
	}()
	return <-runner
}

func (s *goroutine) ScheduleLoop(from int, task func(index int, again func(next int))) Runner {
	runner := make(chan Runner, 1)
	atomic.AddInt32(&s.active, 1)
	s.concurrent.Add(1)
	go func() {
		defer atomic.AddInt32(&s.active, -1)
		defer s.concurrent.Done()
		serial := New()
		runner <- serial.ScheduleLoop(from, task)
		serial.Wait()
	}()
	return <-runner
}

func (s *goroutine) ScheduleFuture(due time.Duration, task func()) Runner {
	cancel := make(cancel)
	atomic.AddInt32(&s.active, 1)
	s.concurrent.Add(1)
	go func() {
		defer atomic.AddInt32(&s.active, -1)
		defer s.concurrent.Done()
		if due > 0 {
			due := time.NewTimer(due)
			select {
			case <-cancel:
				due.Stop()
			case <-due.C:
				task()
			}
		} else {
			select {
			case <-cancel:
				// cancel
			default:
				task()
			}
		}
	}()
	return cancel
}

func (s *goroutine) ScheduleFutureRecursive(due time.Duration, task func(again func(time.Duration))) Runner {
	runner := make(chan Runner, 1)
	atomic.AddInt32(&s.active, 1)
	s.concurrent.Add(1)
	go func() {
		defer atomic.AddInt32(&s.active, -1)
		defer s.concurrent.Done()
		serial := New()
		runner <- serial.ScheduleFutureRecursive(due, task)
		serial.Wait()
	}()
	return <-runner
}

func (s *goroutine) Wait() {
	s.concurrent.Wait()
}

func (s *goroutine) Gosched() {
	runtime.Gosched()
}

func (s *goroutine) IsConcurrent() bool {
	return true
}

func (s *goroutine) Count() int {
	return int(atomic.LoadInt32(&s.active))
}

func (s *goroutine) String() string {
	return fmt.Sprintf("Goroutine{ tasks = %d }", atomic.LoadInt32(&s.active))
}

func (s *goroutine) Concurrent() {
}
