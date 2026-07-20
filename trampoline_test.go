package scheduler_test

import (
	"testing"

	"github.com/reactivego/scheduler"
)

// TestSerialWaitDrainsTasksScheduledAfterPreviousWait: a serial scheduler must
// remain usable after Wait returns; tasks scheduled later are run by the next
// Wait instead of being stranded by a one-shot task loop.
func TestSerialWaitDrainsTasksScheduledAfterPreviousWait(t *testing.T) {
	serial := scheduler.New()
	var got []int
	serial.Schedule(func() { got = append(got, 1) })
	serial.Wait()
	serial.Schedule(func() { got = append(got, 2) })
	serial.Wait()
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("got = %v; want [1 2] (second Wait must run tasks scheduled after the first)", got)
	}
}
