package falcon

import (
	"time"
)

type eventTicker struct {
	t    *Ticker
	done chan struct{}
}

func (this *eventTicker) Run() {
	close(this.done)
	this.t.f()
}

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
			event := &eventTicker{this, ch}
			this.loop.RunInLoop(event)
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
