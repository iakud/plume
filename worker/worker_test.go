package work

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

func TestWork(t *testing.T) {
	work := NewWorker(do, nil)
	work.Join()
	fmt.Printf("work.Join()\n")
}
