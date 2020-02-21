package eventloop

type Pool struct {
	workers []*Worker
	loops   []*EventLoop

	next int
}

func NewPool(numWorkers int, initFunc InitFunc) *Pool {
	var workers []*Worker
	var loops []*EventLoop
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(initFunc)
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
	for _, loop := range this.loops {
		loop.Close()
	}
	for _, worker := range this.workers {
		worker.Join()
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

func (this *Pool) GetAllLoops() []*EventLoop {
	if len(this.loops) == 0 {
		return nil
	}
	loops := make([]*EventLoop, len(this.loops))
	copy(loops, this.loops)
	return loops
}
