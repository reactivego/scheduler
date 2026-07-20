package scheduler

import "context"

// schedulerKey is the context key under which ContextWith attaches a Scheduler.
type schedulerKey struct{}

// ContextWith returns a new context with the provided scheduler attached,
// overriding any scheduler already attached to the parent context.
func ContextWith(parent context.Context, scheduler Scheduler) context.Context {
	return context.WithValue(parent, schedulerKey{}, scheduler)
}

// FromContext returns the Scheduler attached to the context along with true,
// or nil and false when the context has no scheduler attached.
func FromContext(ctx context.Context) (Scheduler, bool) {
	scheduler, ok := ctx.Value(schedulerKey{}).(Scheduler)
	return scheduler, ok
}
