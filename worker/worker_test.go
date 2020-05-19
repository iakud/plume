package worker

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func sleep(ctx context.Context) {
	fmt.Println("sleep second")
	time.Sleep(time.Second)
}

func TestWorker(t *testing.T) {
	worker := NewWorker(sleep)
	worker.Wait()
	fmt.Println("sleep done")
}
