package eventloop

import (
	"log"
	"runtime"
)

type EventLoopWorker struct {
	loop *EventLoop
}

func NewEventLoopWorker(initFunc func(*EventLoop)) *EventLoopWorker {
	worker := &EventLoopWorker{
		loop: NewEventLoop(),
	}
	go worker.workerRoutine(initFunc)
	return worker
}

func (this *EventLoopWorker) GetLoop() *EventLoop {
	return this.loop
}

func (this *EventLoopWorker) workerRoutine(initFunc func(*EventLoop)) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("eventloop: panic worker: %v\n%s", err, buf)
		}
	}()
	if initFunc != nil {
		initFunc(this.loop)
	}
	this.loop.Loop()
}
