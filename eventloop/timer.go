package eventloop

import (
	"time"
)

type Timer struct {
	*time.Timer
}

func newTimer(loop *EventLoop, d time.Duration, f func()) *Timer {
	t := time.AfterFunc(d, func() {
		loop.RunInLoop(f)
	})
	timer := &Timer{t}
	return timer
}
