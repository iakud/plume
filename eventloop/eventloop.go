package eventloop

import (
	"log"
	"runtime"
	"sync"
	"time"
)

type EventLoop struct {
	Context interface{}

	mutex    sync.Mutex
	cond     *sync.Cond
	functors []func()
	closed   bool
}

func NewEventLoop() *EventLoop {
	loop := &EventLoop{}
	loop.cond = sync.NewCond(&loop.mutex)
	return loop
}

func (this *EventLoop) Loop() {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("eventloop: panic loop: %v\n%s", err, buf)
		}
	}()
	var closed bool
	for !closed {
		var functors []func()
		this.mutex.Lock()
		for !this.closed && len(this.functors) == 0 {
			this.cond.Wait()
		}
		functors, this.functors = this.functors, nil // swap
		closed = this.closed
		this.mutex.Unlock()

		for _, functor := range functors {
			functor()
		}
	}
}

func (this *EventLoop) RunInLoop(functor func()) {
	this.mutex.Lock()
	this.functors = append(this.functors, functor)
	this.mutex.Unlock()

	this.cond.Signal()
}

func (loop *EventLoop) Func(functor func()) func() {
	return func() { loop.RunInLoop(functor) }
}

func (this *EventLoop) RunAfter(d time.Duration, f func()) *Timer {
	return newTimer(this, d, f)
}

func (this *EventLoop) RunEvery(d time.Duration, f func()) *Ticker {
	return newTicker(this, d, f)
}

func (this *EventLoop) Close() {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.closed = true
	this.mutex.Unlock()

	this.cond.Signal()
}
