package eventloop

import (
	"time"
)

type Timer struct {
	t *time.Timer
}

func newTimer(loop *EventLoop, d time.Duration, f func()) *Timer {
	timer := &Timer{
		t: time.AfterFunc(d, func() {
			loop.RunInLoop(f)
		}),
	}
	return timer
}

func (this *Timer) Reset(d time.Duration) bool {
	return this.t.Reset(d)
}

func (this *Timer) Stop() bool {
	return this.t.Stop()
}
