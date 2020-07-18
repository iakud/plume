package work

import (
	"context"
	"runtime"
)

type Option interface {
	apply(*WorkerPool)
}

type optionFunc func(*WorkerPool)

func (f optionFunc) apply(pool *WorkerPool) {
	f(pool)
}

var defaultNumWorker = runtime.NumCPU()

// num of workers
func NumWorker(numWorker int) Option {
	return optionFunc(func(pool *WorkerPool) {
		pool.numWorker = numWorker
	})
}

// work proxy
func WorkProxy(workProxy func(ctx context.Context, handler WorkHandler)) Option {
	return optionFunc(func(pool *WorkerPool) {
		if pool.workProxy != nil {
			panic("work: work proxy was already set and may not be reset.")
		}
		pool.workProxy = workProxy
	})
}
