package worker

import (
	"context"
	"log"
	"runtime"
)

type WorkerFunc func(context.Context)

type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorker(workerFunc WorkerFunc) *Worker {
	return NewWorkerWithContext(context.Background(), workerFunc)
}

func NewWorkerWithContext(workerCtx context.Context, workerFunc WorkerFunc) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	worker := &Worker{
		ctx:    ctx,
		cancel: cancel,
	}
	go worker.runner(workerCtx, workerFunc)
	return worker
}

func (this *Worker) Wait() {
	<-this.ctx.Done()
}

func (this *Worker) runner(workerCtx context.Context, workerFunc WorkerFunc) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("worker: panic runner: %v\n%s", err, buf)
		}
		this.cancel()
	}()
	if workerFunc != nil {
		workerFunc(workerCtx)
	}
}
