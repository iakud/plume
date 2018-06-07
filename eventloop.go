package falcon

import (
	"sync"
)

type Event interface {
	Run()
}

type EventLoop struct {
	events []Event
	mutex  sync.Mutex
	cond   *sync.Cond
	closed bool
}

func NewEventLoop() *EventLoop {
	loop := &EventLoop{}
	loop.cond = sync.NewCond(&loop.mutex)
	return loop
}

func (this *EventLoop) Loop() {
	var closed bool = this.IsClosed()
	var events []Event
	for !closed {
		this.mutex.Lock()
		for !this.closed && len(this.events) == 0 {
			this.cond.Wait()
		}
		events = this.events
		this.events = nil
		closed = this.closed
		this.mutex.Unlock()

		for _, event := range events {
			event.Run()
		}
	}
}

func (this *EventLoop) IsClosed() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.closed
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

func (this *EventLoop) RunInLoop(event Event) {
	this.mutex.Lock()
	this.events = append(this.events, event)
	this.mutex.Unlock()

	this.cond.Signal()
}
