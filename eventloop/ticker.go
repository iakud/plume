package eventloop

import (
	"context"
	"time"
)

type Ticker struct {
	t      *time.Ticker
	cancel context.CancelFunc
}

func newTicker(loop *EventLoop, d time.Duration, f func()) *Ticker {
	ctx, cancel := context.WithCancel(context.Background())
	t := time.NewTicker(d)
	ticker := &Ticker{
		t:      t,
		cancel: cancel,
	}
	go func() {
		for {
			select {
			case <-t.C:
				loop.RunInLoop(f)
			case <-ctx.Done():
				return
			}
		}
	}()
	return ticker
}

func (this *Ticker) Stop() {
	this.t.Stop()
	this.cancel()
}
