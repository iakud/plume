package eventloop

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func printLoop(loop *EventLoop) {
	name := loop.Userdata.(string)
	fmt.Printf("print: %s\n", name)
}

var loopId int32 = 0

func onInit(loop *EventLoop) {
	id := atomic.AddInt32(&loopId, 1)
	name := fmt.Sprintf("Loop%d", id)
	loop.Userdata = name
	fmt.Printf("init: %s\n", name)
}

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(3, onInit)
	time.Sleep(time.Millisecond * 100)
	for i := 0; i < 3; i++ {
		nextLoop := pool.GetNextLoop()
		nextLoop.RunInLoop(func() {
			printLoop(nextLoop)
		})
	}
	time.Sleep(time.Millisecond * 100)
}
