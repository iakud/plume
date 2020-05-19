package worker

import (
	"context"
	"errors"
	"sync"
)

var ErrPoolClosed = errors.New("worker: Pool closed")

type PoolHandler interface {
	WorkerContext(ctx context.Context) context.Context
	WorkerExit(ctx context.Context)
}

type TaskFunc func(ctx context.Context)

type Pool struct {
	handler PoolHandler
	maxSize int
	workers []*Worker

	mutex    sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	closed   bool
	queue    []func(context.Context)
}

// maxSize: queue max size, <= 0, Unlimited
func NewPool(numWorkers int, maxSize int, handler PoolHandler) *Pool {
	pool := &Pool{
		maxSize: maxSize,
		handler: handler,
	}
	ctx := context.Background()
	var workers []*Worker
	for i := 0; i < numWorkers; i++ {
		worker := NewWorkerWithContext(ctx, pool.runWorker)
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
	// wait workers exit
	for _, worker := range workers {
		worker.Wait()
	}
}

func (this *Pool) Run(task TaskFunc) error {
	this.mutex.Lock()
	for this.maxSize > 0 && len(this.queue) >= this.maxSize {
		this.notFull.Wait() // wait not full
	}
	if this.closed {
		this.mutex.Unlock() // unlock
		return ErrPoolClosed
	}
	this.queue = append(this.queue, task)
	this.notEmpty.Signal()
	this.mutex.Unlock()
	return nil
}

func (this *Pool) take() (TaskFunc, bool) {
	this.mutex.Lock() // lock
	for !this.closed && len(this.queue) == 0 {
		this.notEmpty.Wait() // wait not empty
	}
	if len(this.queue) > 0 {
		task := this.queue[0]
		this.queue = this.queue[1:]
		this.notFull.Signal() // not full
		this.mutex.Unlock()   // unlock
		return task, true
	}
	this.mutex.Unlock() // unlock
	return nil, false
}

func (this *Pool) runWorker(ctx context.Context) {
	if handler := this.handler; handler != nil {
		ctx = handler.WorkerContext(ctx)
		if ctx == nil {
			panic("WorkerContext returned a nil context")
		}
		defer handler.WorkerExit(ctx)
	}

	for {
		task, ok := this.take()
		if !ok {
			return
		}
		if task != nil {
			task(ctx)
		}
	}
}
