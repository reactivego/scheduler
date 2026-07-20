package scheduler_test

import (
	"context"
	"testing"

	"github.com/reactivego/scheduler"
)

func TestContextWith(t *testing.T) {
	ctx := context.Background()

	if s, ok := scheduler.FromContext(ctx); ok || s != nil {
		t.Fatalf("expected no scheduler in background context, got %v", s)
	}

	serial := scheduler.New()
	ctx = scheduler.ContextWith(ctx, serial)
	if s, ok := scheduler.FromContext(ctx); !ok || s != serial {
		t.Fatalf("expected serial scheduler from context, got %v (ok=%v)", s, ok)
	}

	ctx = scheduler.ContextWith(ctx, scheduler.Goroutine)
	if s, ok := scheduler.FromContext(ctx); !ok || s != scheduler.Scheduler(scheduler.Goroutine) {
		t.Fatalf("expected goroutine scheduler to override serial, got %v (ok=%v)", s, ok)
	}
}
