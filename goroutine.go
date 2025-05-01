package scheduler

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type goroutine struct {
	concurrent sync.WaitGroup
	active     atomic.Int32
}

func (s *goroutine) Concurrent() {
}

func (s *goroutine) Now() time.Time {
	return time.Now()
}

func (s *goroutine) Since(t time.Time) time.Duration {
	return s.Now().Sub(t)
}

func (s *goroutine) Schedule(task func()) Runner {
	cancel := make(cancel)
	s.active.Add(1)
	s.concurrent.Add(1)
	go func() {
		defer s.active.Add(-1)
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
	s.active.Add(1)
	s.concurrent.Add(1)
	go func() {
		defer s.active.Add(-1)
		defer s.concurrent.Done()
		serial := &trampoline{}
		runner <- serial.ScheduleRecursive(task)
		serial.Run()
	}()
	return <-runner
}

func (s *goroutine) ScheduleLoop(from int, task func(index int, again func(next int))) Runner {
	runner := make(chan Runner, 1)
	s.active.Add(1)
	s.concurrent.Add(1)
	go func() {
		defer s.active.Add(-1)
		defer s.concurrent.Done()
		serial := &trampoline{}
		runner <- serial.ScheduleLoop(from, task)
		serial.Run()
	}()
	return <-runner
}

func (s *goroutine) ScheduleFuture(due time.Duration, task func()) Runner {
	cancel := make(cancel)
	s.active.Add(1)
	s.concurrent.Add(1)
	go func() {
		defer s.active.Add(-1)
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
	s.active.Add(1)
	s.concurrent.Add(1)
	go func() {
		defer s.active.Add(-1)
		defer s.concurrent.Done()
		serial := &trampoline{}
		runner <- serial.ScheduleFutureRecursive(due, task)
		serial.Run()
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
	return int(s.active.Load())
}

func (s *goroutine) String() string {
	return fmt.Sprintf("Goroutine{ tasks = %d }", s.active.Load())
}
