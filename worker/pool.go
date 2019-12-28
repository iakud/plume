package worker

import (
	"context"
	"log"
	"runtime"
	"sync"
)

const PoolSizeInfinite = 0

type poolContextKey struct {
}

func (this *poolContextKey) String() string { return "worker context value worker-pool" }

var PoolContextKey = &poolContextKey{}

type PoolWorker interface {
	// WorkerContext修改用于Worker的Context。
	// 提供的Context有一个PoolContextKey值。
	WorkerContext(context.Context, *Worker) context.Context
	WorkerExit(context.Context, *Worker)
}

type Pool struct {
	workers []*Worker
	maxSize int

	ctx        context.Context
	poolWorker PoolWorker

	mutex    sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	closed   bool
	queue    []func(context.Context)
}

func NewPool(numWorkers int, maxSize int, poolWorker PoolWorker) *Pool {
	pool := &Pool{
		maxSize: maxSize,

		poolWorker: poolWorker,
	}
	pool.ctx = context.WithValue(context.Background(), PoolContextKey, pool)
	var workers []*Worker
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(pool.workerRoutine)
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

func (this *Pool) Run(task func(context.Context)) {
	this.mutex.Lock()
	for this.maxSize > 0 && len(this.queue) >= this.maxSize {
		this.notFull.Wait()
	}
	this.queue = append(this.queue, task)
	this.mutex.Unlock()
	this.notEmpty.Signal()
}

func (this *Pool) workerRoutine(worker *Worker) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("worker: panic worker: %v\n%s", err, buf)
		}
	}()
	ctx := this.ctx
	if poolWorker := this.poolWorker; poolWorker != nil {
		ctx = poolWorker.WorkerContext(ctx, worker)
		if ctx == nil {
			panic("WorkerContext returned a nil context")
		}
		defer poolWorker.WorkerExit(ctx, worker)
	}
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
