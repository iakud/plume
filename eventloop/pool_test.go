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
	fmt.Printf("%s: run in loop\n", name)
}

type loopPool struct {
	loopId int32
}

func (this *loopPool) LoopStart(loop *EventLoop) {
	id := atomic.AddInt32(&this.loopId, 1)
	name := fmt.Sprintf("Loop%d", id)
	loop.Context = name
	fmt.Printf("%s: start\n", name)
}

func (this *loopPool) LoopStop(loop *EventLoop) {
	name, ok := loop.Context.(string)
	if !ok {
		return
	}
	fmt.Printf("%s: stop\n", name)
}

func TestPool(t *testing.T) {
	pool := NewPool(3, &loopPool{})
	time.Sleep(time.Millisecond * 100)
	for i := 0; i < 3; i++ {
		nextLoop := pool.GetNextLoop()
		nextLoop.RunInLoop(func() {
			printName(nextLoop)
		})
	}
	pool.Close()
}
