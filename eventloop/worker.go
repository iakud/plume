package eventloop

import (
	"log"
	"runtime"
	"sync"
)

type InitFunc func(loop *EventLoop)

type Worker struct {
	initFunc InitFunc
	loop     *EventLoop

	done sync.WaitGroup
}

func NewWorker(initFunc InitFunc) *Worker {
	worker := &Worker{
		loop:     NewEventLoop(),
		initFunc: initFunc,
	}
	worker.done.Add(1)
	go func() {
		defer worker.done.Done()
		worker.worker()
	}()
	return worker
}

func (this *Worker) GetLoop() *EventLoop {
	return this.loop
}

func (this *Worker) Join() {
	this.done.Wait()
}

func (this *Worker) worker() {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("eventloop: panic worker: %v\n%s", err, buf)
		}
	}()
	if initFunc := this.initFunc; initFunc != nil {
		initFunc(this.loop)
	}
	this.loop.Loop()
}
