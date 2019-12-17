package eventloop

import (
	"math/rand"
	"time"
)

type EventLoopWorkerPool struct {
	workers []*EventLoopWorker
	loops   []*EventLoop

	rand *rand.Rand
	next int
}

func NewEventLoopWorkerPool(numWorkers int, initFunc func(*EventLoop)) *EventLoopWorkerPool {
	pool := &EventLoopWorkerPool{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	var workers []*EventLoopWorker
	var loops []*EventLoop
	for i := 0; i < numWorkers; i++ {
		worker := NewEventLoopWorker(initFunc)
		workers = append(workers, worker)
		loops = append(loops, worker.GetLoop())
	}
	pool.workers = workers
	pool.loops = loops
	return pool
}

func (this *EventLoopWorkerPool) GetNextLoop() *EventLoop {
	if len(this.loops) == 0 {
		return nil
	}
	index := this.next
	this.next++
	if this.next >= len(this.loops) {
		this.next = 0
	}
	return this.loops[index]
}

func (this *EventLoopWorkerPool) GetRandLoop() *EventLoop {
	if len(this.loops) == 0 {
		return nil
	}
	return this.loops[rand.Int()%len(this.loops)]
}

func (this *EventLoopWorkerPool) GetAllLoops() []*EventLoop {
	if len(this.loops) == 0 {
		return nil
	}
	loops := make([]*EventLoop, len(this.loops))
	copy(loops, this.loops)
	return loops
}
