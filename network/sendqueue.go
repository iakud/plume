package network

import (
	"errors"
	"sync"
)

var ErrWriteClosed = errors.New("write closed")

type SendQueue struct {
	queue [][]byte

	mutex  sync.Mutex
	cond   *sync.Cond
	closed bool
}

func NewSendQueue() *SendQueue {
	sendQueue := new(SendQueue)
	sendQueue.cond = sync.NewCond(&sendQueue.mutex)
	return sendQueue
}

func (this *SendQueue) Append(b []byte) error {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return ErrWriteClosed
	}
	this.queue = append(this.queue, b)
	this.mutex.Unlock()

	this.cond.Signal()
	return nil
}

func (this *SendQueue) Get() [][]byte {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for !this.closed && len(this.queue) == 0 {
		this.cond.Wait()
	}
	queue := this.queue
	this.queue = nil
	return queue
}

func (this *SendQueue) Close() {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.closed = true
	this.mutex.Unlock()

	this.cond.Signal()
}
