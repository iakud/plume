package eventloop

import (
	"fmt"
	"testing"
)

func quit(loop *EventLoop) {
	fmt.Printf("loop quit\n")
	loop.Close()
}

func TestWorker(t *testing.T) {
	done := make(chan struct{})
	worker := NewWorker(func(loop *EventLoop) {
		fmt.Printf("loop init\n")
	})
	loop := worker.GetLoop()
	loop.RunInLoop(func() {
		quit(loop)
		close(done)
	})
	<-done
}
