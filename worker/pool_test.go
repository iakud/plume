package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var workerId int32 = 0

type workerNameKey struct {
}

func newWorkerNameContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, workerNameKey{}, name)
}

func fromWorkerNameContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(workerNameKey{}).(string)
	return name, ok
}

func runWorker(ctx context.Context, handler WorkerHandler) {
	id := atomic.AddInt32(&workerId, 1)
	name := fmt.Sprintf("worker%d", id)
	fmt.Printf("%s init\n", name)
	defer func() {
		fmt.Printf("%s exit\n", name)
	}()

	ctx = newWorkerNameContext(ctx, name)
	handler(ctx)
}

func TestPool(t *testing.T) {
	pool := NewPool(3, PoolSizeInfinite, runWorker)
	time.Sleep(time.Second)
	for i := 0; i < 100; i++ {
		buf := fmt.Sprintf("task %d", i)
		pool.Run(func(ctx context.Context) {
			name, ok := fromWorkerNameContext(ctx)
			if !ok {
				return
			}
			fmt.Printf("%s run: %s\n", name, buf)
			time.Sleep(time.Millisecond * 100)
		})
	}
	var wg sync.WaitGroup
	wg.Add(1)
	pool.Run(func(ctx context.Context) {
		wg.Done()
	})
	wg.Wait()
	pool.Close()
}
