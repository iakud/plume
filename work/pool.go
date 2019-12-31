package work

import (
	"context"
	"sync"
)

const PoolSizeInfinite = 0

type poolContextKey struct {
}

func (this *poolContextKey) String() string { return "worker context value worker-pool" }

var PoolContextKey = &poolContextKey{}

type RunnerInterceptor func(ctx context.Context, handler RunnerHandler)

type RunnerHandler func(ctx context.Context)

type Pool struct {
	works   []*Work
	maxSize int

	ctx       context.Context
	runnerInt RunnerInterceptor

	mutex    sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	closed   bool
	queue    []func(context.Context)
}

func NewPool(numWorkers int, maxSize int, runnerInt RunnerInterceptor) *Pool {
	pool := &Pool{
		maxSize: maxSize,

		runnerInt: runnerInt,
	}
	pool.ctx = context.WithValue(context.Background(), PoolContextKey, pool)
	var works []*Work
	for i := 0; i < numWorkers; i++ {
		work := NewWork(pool.workRunner)
		works = append(works, work)
	}
	pool.works = works
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
	works := this.works
	this.works = nil
	this.mutex.Unlock()
	for _, work := range works {
		work.Join()
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

func (this *Pool) runner(ctx context.Context) {
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

func (this *Pool) workRunner() {
	if this.runnerInt == nil {
		this.runner(this.ctx)
		return
	}
	this.runnerInt(this.ctx, this.runner)
}
