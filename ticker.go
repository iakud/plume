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
	go ticker.serve()
	return ticker
}

func (this *Ticker) serve() {
	for {
		select {
		case _ = <-this.t.C:
			this.onTicker()
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

func (this *Ticker) onTicker() {
	ch := make(chan struct{})
	if this.loop == nil {
		close(ch)
		this.f()
	} else {
		this.loop.RunInLoop(func() {
			close(ch)
			this.f()
		})
	}
	select {
	case <-ch:
	case <-this.done:
		return
	}
}
