package worker

import (
	"context"
	"log"
	"runtime"
	"sync"
)

type WorkerFunc func(ctx context.Context)

type Worker struct {
	f  WorkerFunc
	wg sync.WaitGroup
}

func NewWorker(f WorkerFunc) *Worker {
	return NewWorkerWithContext(context.Background(), f)
}

func NewWorkerWithContext(ctx context.Context, f WorkerFunc) *Worker {
	worker := &Worker{
		f: f,
	}
	worker.wg.Add(1)
	go worker.runner(ctx)
	return worker
}

func (this *Worker) Join() {
	this.wg.Wait()
}

func (this *Worker) runner(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("work: panic worker: %v\n%s", err, buf)
		}
		this.wg.Done()
	}()
	this.f(ctx)
}
