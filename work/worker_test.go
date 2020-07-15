package work

import (
	"fmt"
	"testing"
	"time"
)

func sleep() {
	fmt.Println("sleep second")
	time.Sleep(time.Second)
}

func TestWorker(t *testing.T) {
	worker := NewWorker(sleep)
	<-worker.Done()
	fmt.Println("worker done")
}
