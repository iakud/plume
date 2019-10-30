package eventloop

import (
	"context"
	"time"
)

type Ticker struct {
	cancel context.CancelFunc
}

func newTicker(loop *EventLoop, d time.Duration, f func()) *Ticker {
	ctx, cancel := context.WithCancel(context.Background())
	ticker := &Ticker{
		cancel: cancel,
	}
	go func() {
		t := time.NewTicker(d)
		defer t.Stop()
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
	this.cancel()
}
