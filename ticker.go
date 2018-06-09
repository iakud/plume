package falcon

import (
	"time"
)

type Ticker struct {
	loop *EventLoop
	t    *time.Ticker
	f    func()
	done chan struct{}
}

func NewTicker(loop *EventLoop, d time.Duration, f func()) *Ticker {
	t := time.NewTicker(d)
	ticker := &Ticker{
		loop: loop,
		t:    t,
		f:    f,
		done: make(chan struct{}),
	}
	go ticker.ticker()
	return ticker
}

func (this *Ticker) ticker() {
	for {
		select {
		case _ = <-this.t.C:
			ch := make(chan struct{})
			this.loop.RunInLoop(func() {
				close(ch)
				this.f()
			})

			select {
			case <-ch:
			case <-this.done:
				return
			}
		case <-this.done:
			return
		}
	}
}

func (this *Ticker) Stop() {
	this.t.Stop()
	select {
	case <-this.done:
	default:
		close(this.done)
	}
}
