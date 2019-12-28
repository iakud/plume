package worker

import (
	"context"
)

type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc

	workerFunc func(*Worker)
}

func NewWorker(workerFunc func(*Worker)) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	worker := &Worker{
		ctx:    ctx,
		cancel: cancel,

		workerFunc: workerFunc,
	}
	go worker.workerRoutine()
	return worker
}

func (this *Worker) workerRoutine() {
	defer this.cancel()
	this.workerFunc(this)
}

func (this *Worker) Join() {
	<-this.ctx.Done()
}
