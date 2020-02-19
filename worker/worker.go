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

	done sync.WaitGroup
}

func NewWorker(workerFunc WorkerFunc) *Worker {
	return NewWorkerWithContext(context.Background(), workerFunc)
}

func NewWorkerWithContext(ctx context.Context, workerFunc WorkerFunc) *Worker {
	worker := &Worker{
		workerFunc: workerFunc,
	}
	worker.done.Add(1)
	go func() {
		defer worker.done.Done()
		worker.worker(ctx)
	}()
	return worker
}

func (this *Worker) Join() {
	this.done.Wait()
}

func (this *Worker) worker(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("work: panic worker: %v\n%s", err, buf)
		}
	}()
	this.workerFunc(ctx)
}
