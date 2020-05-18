package eventloop

import (
	"fmt"
	"testing"
)

type loopWorker struct {
}

func (this *loopWorker) LoopInit(loop *EventLoop) {
	fmt.Printf("loop init\n")
}

func (this *loopWorker) LoopClose(loop *EventLoop) {
	fmt.Printf("loop close\n")
}

func TestWorker(t *testing.T) {
	worker := NewWorker(&loopWorker{})
	loop := worker.GetLoop()
	loop.RunInLoop(func() {
		fmt.Printf("in loop\n")
	})
	worker.Close()
}
