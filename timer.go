package falcon

import (
	"time"
)

type eventTimer struct {
	t *Timer
}

func (this *eventTimer) Run() {
	this.t.f()
}

type Timer struct {
	loop *EventLoop
	t    *time.Timer
	f    func()
	done chan struct{}
}

func NewTimer(loop *EventLoop, d time.Duration, f func()) *Timer {
	t := time.NewTimer(d)
	timer := &Timer{
		loop: loop,
		t:    t,
		f:    f,
		done: make(chan struct{}),
	}
	go timer.timer()
	return timer
}

func (this *Timer) timer() {
	select {
	case <-this.t.C:
		event := &eventTimer{this}
		this.loop.RunInLoop(event)
	case <-this.done:

	}
}

func (this *Timer) Stop() {
	this.t.Stop()
	select {
	case <-this.done:
	default:
		close(this.done)
	}
}
