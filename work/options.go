package work

import (
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

// worker context
func WorkerCtx(workerCtx WorkerContext) Option {
	return optionFunc(func(pool *WorkerPool) {
		if pool.workerCtx != nil {
			panic("work: worker context was already set and may not be reset.")
		}
		pool.workerCtx = workerCtx
	})
}
