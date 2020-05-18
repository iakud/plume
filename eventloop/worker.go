package eventloop

import (
	"log"
	"runtime"
	"sync"
)

type InitFunc func(loop *EventLoop)

type Worker struct {
	loop     *EventLoop
	initFunc InitFunc

	initWg sync.WaitGroup
	exitWg sync.WaitGroup
}

func NewWorker(initFunc InitFunc) *Worker {
	worker := &Worker{
		loop:     NewEventLoop(),
		initFunc: initFunc,
	}
	worker.initWg.Add(1)
	worker.exitWg.Add(1)
	go worker.runLoop()
	worker.initWg.Wait() // return after initFunc
	return worker
}

func (this *Worker) Close() {
	this.loop.Close()
	this.exitWg.Wait()
}

func (this *Worker) GetLoop() *EventLoop {
	return this.loop
}

func (this *Worker) runLoop() {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("eventloop: panic worker: %v\n%s", err, buf)
		}
		this.exitWg.Done()
	}()
	if initFunc := this.initFunc; initFunc != nil {
		initFunc(this.loop)
	}
	this.initWg.Done()
	this.loop.Loop()
}
