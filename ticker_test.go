package falcon

import (
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	loop := NewEventLoop()
	times := 0
	NewTicker(loop, time.Second, func() {
		times++
		if times == 3 {
			loop.Close()
		}
	})
	loop.Loop()
}
