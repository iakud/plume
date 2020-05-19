package concurrent

import (
	"sync"
)

type CountDownLatch struct {
	mutex sync.Mutex
	cond  *sync.Cond
	count int
}

func NewCountDownLatch(count int) *CountDownLatch {
	countDownLatch := &CountDownLatch{
		count: count,
	}
	countDownLatch.cond = sync.NewCond(&countDownLatch.mutex)
	return countDownLatch
}

func (this *CountDownLatch) Wait() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for this.count > 0 {
		this.cond.Wait()
	}
}

func (this *CountDownLatch) CountDown() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.count--
	if this.count == 0 {
		this.cond.Broadcast()
	}
}
