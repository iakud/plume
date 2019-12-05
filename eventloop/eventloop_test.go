package eventloop

import (
	"fmt"
	"testing"
)

func TestEventLoop(t *testing.T) {
	loop := NewEventLoop()
	loop.RunInLoop(func() {
		fmt.Println("close in loop")
		loop.Close()
	})
	loop.Loop()
}
