package work

import (
	"context"
	"log"
	"runtime"
	"sync"
)

type Handler interface {
	WorkerContext(ctx context.Context, worker *Worker) context.Context
	Exit(ctx context.Context)
}

type WorkFunc func(ctx context.Context)

type Worker struct {
	workFunc WorkFunc
	handler  Handler

	done sync.WaitGroup
}

func NewWorker(workFunc WorkFunc, handler Handler) *Worker {
	return NewWorkerContext(context.Background(), workFunc, handler)
}

func NewWorkerContext(ctx context.Context, workFunc WorkFunc, handler Handler) *Worker {
	worker := &Worker{
		workFunc: workFunc,
		handler:  handler,
	}
	worker.done.Add(1)
	go func() {
		defer worker.done.Done()
		worker.runner(ctx)
	}()
	return worker
}

func (this *Worker) Join() {
	this.done.Wait()
}

func (this *Worker) runner(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("work: panic runner: %v\n%s", err, buf)
		}
	}()
	if handler := this.handler; handler != nil {
		ctx = handler.WorkerContext(ctx, this)
		if ctx == nil {
			panic("WorkerContext returned a nil context")
		}
		defer handler.Exit(ctx)
	}
	this.workFunc(ctx)
}
