package work

import (
	"context"
	"runtime"
)

type options struct {
	numWorker int
	workProxy func(ctx context.Context, handler WorkHandler)
}

var defaultOptions = options{
	numWorker: runtime.NumCPU(),
}

type Option func(*options)

// num of workers
func NumWorker(numWorker int) Option {
	return func(opts *options) {
		opts.numWorker = numWorker
	}
}

// work proxy
func WorkProxy(workProxy func(ctx context.Context, handler WorkHandler)) Option {
	return func(opts *options) {
		if opts.workProxy != nil {
			panic("work: work proxy was already set and may not be reset.")
		}
		opts.workProxy = workProxy
	}
}
