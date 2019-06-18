package eventloop

import (
	"time"
)

type Ticker struct {
	loop *EventLoop
	t    *time.Ticker
	f    func(time.Time)
	done chan struct{}
}

func newTicker(loop *EventLoop, d time.Duration, f func(time.Time)) *Ticker {
	t := time.NewTicker(d)
	ticker := &Ticker{
		loop: loop,
		t:    t,
		f:    f,
		done: make(chan struct{}),
	}
	go ticker.receiveTime()
	return ticker
}

func (this *Ticker) receiveTime() {
	for {
		select {
		case now := <-this.t.C:
			this.loop.RunInLoop(func() {
				this.f(now)
			})
		case <-this.done:
			return
		}
	}
}

func (this *Ticker) Stop() {
	select {
	case <-this.done:
	default:
		close(this.done)
		this.t.Stop()
	}
}
