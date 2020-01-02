package worker

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func do(ctx context.Context) {
	time.Sleep(time.Second)
	fmt.Printf("wait()\n")
}

func TestWorker(t *testing.T) {
	worker := NewWorker(do, nil)
	worker.Join()
	fmt.Printf("worker.Join()\n")
}
