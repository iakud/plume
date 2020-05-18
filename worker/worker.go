package worker

import (
	"context"
	"log"
	"runtime"
	"sync"
)

type WorkerFunc func(ctx context.Context)

type Worker struct {
	workerFunc WorkerFunc
	exitWg     sync.WaitGroup
}

func NewWorker(f WorkerFunc) *Worker {
	return NewWorkerWithContext(context.Background(), f)
}

func NewWorkerWithContext(ctx context.Context, workerFunc WorkerFunc) *Worker {
	worker := &Worker{
		workerFunc: workerFunc,
	}
	worker.exitWg.Add(1)
	go worker.runner(ctx)
	return worker
}

func (this *Worker) Join() {
	this.exitWg.Wait()
}

func (this *Worker) runner(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("worker: panic runner: %v\n%s", err, buf)
		}
		this.exitWg.Done()
	}()
	if workerFunc := this.workerFunc; workerFunc != nil {
		workerFunc(ctx)
	}
}
