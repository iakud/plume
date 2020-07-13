package worker

import (
	"context"
	"errors"
	//	"sync"
)

var ErrPoolClosed = errors.New("worker: Pool closed")

type PoolHandler interface {
	WorkerContext(ctx context.Context) context.Context
	WorkerExit(ctx context.Context)
}

type TaskFunc func(ctx context.Context)

type Pool struct {
	maxSize int
	handler PoolHandler
	workers []*Worker

	taskCh chan TaskFunc
	/*
		mutex    sync.Mutex
		notFull  *sync.Cond
		notEmpty *sync.Cond
		closed   bool
		queue    []func(context.Context)
	*/
}

func NewPool(numWorkers int, maxSize int, handler PoolHandler) *Pool {
	pool := &Pool{
		maxSize: maxSize,
		handler: handler,

		taskCh: make(chan TaskFunc, maxSize),
	}
	ctx := context.Background()
	var workers []*Worker
	for i := 0; i < numWorkers; i++ {
		worker := NewWorkerWithContext(ctx, pool.runWorker)
		workers = append(workers, worker)
	}
	pool.workers = workers
	return pool
}

func (this *Pool) Close() {
	close(this.taskCh)
	for _, worker := range this.workers {
		worker.Wait()
	}
	this.workers = nil
}

func (this *Pool) Run(ctx context.Context, task TaskFunc) error {
	if task == nil {
		return nil
	}
	select {
	case this.taskCh <- task:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (this *Pool) runWorker(ctx context.Context) {
	if handler := this.handler; handler != nil {
		ctx = handler.WorkerContext(ctx)
		if ctx == nil {
			panic("WorkerContext returned a nil context")
		}
		defer handler.WorkerExit(ctx)
	}

	for task := range this.taskCh {
		task(ctx)
	}
}
