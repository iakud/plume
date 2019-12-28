package worker

import (
	"fmt"
	"testing"
	"time"
)

func waitSecond(worker *Worker) {
	time.Sleep(time.Second)
	fmt.Printf("waitSecond()\n")
}

func TestWorker(t *testing.T) {
	worker := NewWorker(waitSecond)
	worker.Join()
	fmt.Printf("worker.Join()\n")
}
