package falcon

import (
	"testing"
)

func TestEventLoop(t *testing.T) {
	loop := NewEventLoop()
	loop.Close()
	loop.Loop()
}
