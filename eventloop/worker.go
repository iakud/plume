package eventloop

import (
	"log"
	"runtime"
)

type InitFunc func(loop *EventLoop)

type Worker struct {
	initFunc InitFunc
	loop     *EventLoop
}

func NewWorker(initFunc InitFunc) *Worker {
	worker := &Worker{
		loop:     NewEventLoop(),
		initFunc: initFunc,
	}
	go worker.runner()
	return worker
}

func (this *Worker) GetLoop() *EventLoop {
	return this.loop
}

func (this *Worker) runner() {
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
