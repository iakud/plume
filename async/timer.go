package falcon

import (
	"time"
)

type Timer struct {
	loop *EventLoop
	t    *time.Timer
	f    func(time.Time)
	done chan struct{}
}

func newTimer(loop *EventLoop, d time.Duration, f func(time.Time)) *Timer {
	t := time.NewTimer(d)
	timer := &Timer{
		loop: loop,
		t:    t,
		f:    f,
		done: make(chan struct{}),
	}
	go timer.receiveTime()
	return timer
}

func (this *Timer) receiveTime() {
	select {
	case now := <-this.t.C:
		this.loop.RunInLoop(func() {
			this.f(now)
		})
	case <-this.done:
	}
}

func (this *Timer) Stop() bool {
	if this.t.Stop() {
		close(this.done)
		return true
	}
	return false
}
