package eventloop

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func printName(loop *EventLoop) {
	name, ok := loop.Context.(string)
	if !ok {
		return
	}
	fmt.Printf("loop name: %s\n", name)
}

var loopId int32 = 0

func onLoopInit(loop *EventLoop) {
	id := atomic.AddInt32(&loopId, 1)
	name := fmt.Sprintf("Loop%d", id)
	loop.Context = name
	fmt.Printf("init: %s\n", name)
}

func TestPool(t *testing.T) {
	pool := NewPool(3, onLoopInit)
	time.Sleep(time.Millisecond * 100)
	for i := 0; i < 3; i++ {
		nextLoop := pool.GetNextLoop()
		nextLoop.RunInLoop(func() {
			printName(nextLoop)
		})
	}
	pool.Close()
}
