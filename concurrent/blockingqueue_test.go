package concurrent

import (
	"testing"
)

const MSG_COUNT int = 5000000
const GO_NUM int = 2

func bqPut(bq *BlockingQueue) {
	for i := 0; i < MSG_COUNT; i++ {
		bq.Put(i)
	}
}

func TestBlockQueue(t *testing.T) {
	bq := NewBlockingQueue()
	for i := 0; i < GO_NUM; i++ {
		go bqPut(bq)
	}
	times := 0
	for {
		_ = bq.Take()
		times++
		if times == MSG_COUNT*GO_NUM {
			return
		}
	}
}

func sendToChan(ch chan interface{}) {
	for i := 0; i < MSG_COUNT; i++ {
		ch <- i
	}
}

func TestChan(t *testing.T) {
	ch := make(chan interface{}, 1024)
	for i := 0; i < GO_NUM; i++ {
		go sendToChan(ch)
	}
	times := 0
	for {
		_ = <-ch
		times++
		if times == MSG_COUNT*GO_NUM {
			return
		}
	}
}
