package concurrent

import (
	"sync"
)

type BlockingQueue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	queue []interface{}
}

func NewBlockingQueue() *BlockingQueue {
	bq := &BlockingQueue{}
	bq.cond = sync.NewCond(&bq.mu)
	return bq
}

func (this *BlockingQueue) Put(v interface{}) {
	this.mu.Lock()
	this.queue = append(this.queue, v)
	this.mu.Unlock()
	this.cond.Signal()
}

func (this *BlockingQueue) Take() interface{} {
	this.mu.Lock()
	for len(this.queue) == 0 {
		this.cond.Wait()
	}
	v := this.queue[0]
	this.queue = this.queue[1:]
	this.mu.Unlock()
	return v
}

func (this *BlockingQueue) Len() int {
	return len(this.queue)
}
