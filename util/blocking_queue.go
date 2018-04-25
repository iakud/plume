package util

import (
	"container/list"
	"sync"
)

type BlockingQueue struct {
	mu     sync.Mutex
	cond   *sync.Cond
	list   *list.List
	inWait int
}

func NewBlockingQueue() *BlockingQueue {
	bq := &BlockingQueue{}
	bq.cond = sync.NewCond(&bq.mu)
	bq.list = list.New()
	return bq
}

func (this *BlockingQueue) Put(v interface{}) {
	this.mu.Lock()
	this.list.PushBack(v)
	inWait := this.inWait
	this.mu.Unlock()
	if inWait > 0 {
		this.cond.Signal() // 有Wait才Signal
	}
}

func (this *BlockingQueue) Take() interface{} {
	this.mu.Lock()
	for this.list.Len() == 0 {
		this.inWait++
		this.cond.Wait()
		this.inWait--
	}
	v := this.list.Remove(this.list.Front())
	this.mu.Unlock()
	return v
}
