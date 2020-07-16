package work

import (
	"context"
)

type TaskFunc func(ctx context.Context)

type WorkerContext interface {
	WorkContext(context.Context) context.Context
	WorkExit(context.Context)
}

type WorkerPool struct {
	taskCh  chan TaskFunc
	workers []*Worker

	numWorker int
	workerCtx WorkerContext
}

func NewWorkerPool(size int, opts ...Option) *WorkerPool {
	pool := &WorkerPool{
		taskCh:    make(chan TaskFunc, size),
		numWorker: defaultNumWorker,
	}
	// options apply
	for _, opt := range opts {
		opt.apply(pool)
	}
	// workers run
	var workers []*Worker
	for i := 0; i < pool.numWorker; i++ {
		worker := NewWorker(pool.process)
		workers = append(workers, worker)
	}
	pool.workers = workers
	return pool
}

func (pool *WorkerPool) Close() {
	close(pool.taskCh)
	// wait workers done
	for _, worker := range pool.workers {
		<-worker.Done()
	}
	pool.workers = nil
}

func (pool *WorkerPool) Run(task TaskFunc) {
	pool.taskCh <- task
}

func (pool *WorkerPool) RunContext(ctx context.Context, task TaskFunc) error {
	select {
	case pool.taskCh <- task:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (pool *WorkerPool) TryRun(task TaskFunc) bool {
	select {
	case pool.taskCh <- task:
	default:
		return false
	}
	return true
}

func (pool *WorkerPool) process() {
	ctx := context.Background()

	if workerCtx := pool.workerCtx; workerCtx != nil {
		ctx = workerCtx.WorkContext(ctx)
		if ctx == nil {
			panic("work: WorkContext returned a nil context")
		}
		defer workerCtx.WorkExit(ctx)
	}
	for task := range pool.taskCh {
		if task != nil {
			task(ctx)
		}
	}
}
