package eventloop

import (
	"fmt"
	"testing"
)

func quit(loop *EventLoop) {
	name := loop.Context.(string)
	fmt.Printf("%s close\n", name)
	loop.Close()
}

func workerInit(loop *EventLoop) {
	name := "red"
	loop.Context = name
}

func TestEventLoopWorker(t *testing.T) {
	worker := NewEventLoopWorker(workerInit)
	loop := worker.GetLoop()
	done := make(chan struct{})
	loop.RunInLoop(func() {
		quit(loop)
		close(done)
	})
	<-done
}
