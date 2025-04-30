package scheduler

import "time"

type futuretask struct {
	at     time.Time
	run    func()
	cancel cancel
}
