package eventloop

import (
	"math/rand"
)

type WorkerPool struct {
	workers []*Worker
	loops   []*EventLoop

	rand *rand.Rand
	next int
}

func NewWorkerPool(numWorkers int, initFunc InitFunc) *WorkerPool {
	var workers []*Worker
	var loops []*EventLoop
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(initFunc)
		workers = append(workers, worker)
		loops = append(loops, worker.GetLoop())
	}
	pool := &WorkerPool{
		workers: workers,
		loops:   loops,
	}
	return pool
}

func (this *WorkerPool) GetNextLoop() *EventLoop {
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

func (this *WorkerPool) GetAllLoops() []*EventLoop {
	if len(this.loops) == 0 {
		return nil
	}
	loops := make([]*EventLoop, len(this.loops))
	copy(loops, this.loops)
	return loops
}
