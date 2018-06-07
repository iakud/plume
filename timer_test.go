package falcon

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	loop := NewEventLoop()
	NewTimer(loop, time.Second, func() {
		loop.Close()
	})
	loop.Loop()
}
