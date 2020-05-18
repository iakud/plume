package eventloop

import (
	"log"
	"runtime"
	"sync"
)

type LoopHandler interface {
	LoopStart(loop *EventLoop)
	LoopStop(loop *EventLoop)
}

type InitFunc func(loop *EventLoop)

type Worker struct {
	loop    *EventLoop
	handler LoopHandler

	initWg sync.WaitGroup
	exitWg sync.WaitGroup
}

func NewWorker(handler LoopHandler) *Worker {
	worker := &Worker{
		loop:    NewEventLoop(),
		handler: handler,
	}
	worker.initWg.Add(1)
	worker.exitWg.Add(1)
	go worker.runLoop()
	worker.initWg.Wait()
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
			log.Printf("eventloop: panic runLoop: %v\n%s", err, buf)
		}
		this.exitWg.Done()
	}()
	if handler := this.handler; handler != nil {
		handler.LoopStart(this.loop)
		defer handler.LoopStop(this.loop)
	}
	this.initWg.Done()
	this.loop.Loop()
}
