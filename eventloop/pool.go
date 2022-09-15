package eventloop

type Pool struct {
	workers []*Worker
	loops   []*EventLoop

	next int
}

func NewPool(numWorkers int, handler LoopHandler) *Pool {
	var workers []*Worker
	var loops []*EventLoop
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(handler)
		workers = append(workers, worker)
		loops = append(loops, worker.GetLoop())
	}
	pool := &Pool{
		workers: workers,
		loops:   loops,
	}
	return pool
}

func (this *Pool) Close() {
	for _, worker := range this.workers {
		worker.Close()
	}
}

func (this *Pool) GetNextLoop() *EventLoop {
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

func (this *Pool) GetLoopForHash(hashCode int) *EventLoop {
	if len(this.loops) == 0 {
		return nil
	}
	index := hashCode % len(this.loops)
	return this.loops[index]
}

func (this *Pool) GetAllLoops() []*EventLoop {
	if len(this.loops) == 0 {
		return nil
	}
	loops := make([]*EventLoop, len(this.loops))
	copy(loops, this.loops)
	return loops
}
