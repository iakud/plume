package eventloop

import (
	"fmt"
	"testing"
)

type loopWorker struct {
}

func (this *loopWorker) LoopStart(loop *EventLoop) {
	fmt.Printf("loop start\n")
}

func (this *loopWorker) LoopStop(loop *EventLoop) {
	fmt.Printf("loop stop\n")
}

func TestWorker(t *testing.T) {
	worker := NewWorker(&loopWorker{})
	loop := worker.GetLoop()
	loop.RunInLoop(func() {
		fmt.Printf("in loop\n")
	})
	worker.Close()
}
