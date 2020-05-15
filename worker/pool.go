package worker

import (
	"context"
	"sync"
)

type WorkerHandler func(ctx context.Context)
type WorkerInterceptor func(ctx context.Context, handler WorkerHandler)

const PoolSizeInfinite = 0

type TaskFunc func(ctx context.Context)

type Pool struct {
	workers     []*Worker
	maxSize     int
	interceptor WorkerInterceptor

	mutex    sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	closed   bool
	queue    []func(context.Context)
}

func NewPool(numWorkers int, maxSize int, interceptor WorkerInterceptor) *Pool {
	pool := &Pool{
		maxSize:     maxSize,
		interceptor: interceptor,
	}
	ctx := context.Background()
	var workers []*Worker
	for i := 0; i < numWorkers; i++ {
		worker := NewWorkerWithContext(ctx, pool.runner)
		workers = append(workers, worker)
	}
	pool.workers = workers
	pool.notFull = sync.NewCond(&pool.mutex)
	pool.notEmpty = sync.NewCond(&pool.mutex)
	return pool
}

func (this *Pool) Close() {
	this.mutex.Lock()
	if this.closed {
		return
	}
	this.closed = true
	this.notEmpty.Broadcast()
	workers := this.workers
	this.workers = nil
	this.mutex.Unlock()
	for _, worker := range workers {
		worker.Join()
	}
}

func (this *Pool) Run(task TaskFunc) {
	this.mutex.Lock()
	for this.maxSize > 0 && len(this.queue) >= this.maxSize {
		this.notFull.Wait()
	}
	this.queue = append(this.queue, task)
	this.mutex.Unlock()
	this.notEmpty.Signal()
}

func (this *Pool) runner(ctx context.Context) {
	if this.interceptor == nil {
		this.runWorker(ctx)
	}
	this.interceptor(ctx, this.runWorker)
}

func (this *Pool) runWorker(ctx context.Context) {
	var closed bool
	for !closed {
		this.mutex.Lock()
		for !this.closed && len(this.queue) == 0 {
			this.notEmpty.Wait()
		}
		var task func(context.Context)
		var notFull bool
		if len(this.queue) > 0 {
			task = this.queue[0]
			this.queue = this.queue[1:]
			notFull = this.maxSize > 0
		}
		closed = this.closed
		this.mutex.Unlock()
		if notFull {
			this.notFull.Signal()
		}
		if task != nil {
			task(ctx)
		}
	}
}
