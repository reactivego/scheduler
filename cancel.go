package scheduler

type cancel chan struct{}

func (c cancel) Cancel() {
	close(c)
}
