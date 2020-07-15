package work

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

var workerId int32 = 0

type workerNameKey struct{}

func newWorkerNameContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, workerNameKey{}, name)
}

func fromWorkerNameContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(workerNameKey{}).(string)
	return name, ok
}

func namedWorker(ctx context.Context, handler WorkerHandler) {
	name := fmt.Sprintf("worker%d", atomic.AddInt32(&workerId, 1))
	ctx = newWorkerNameContext(ctx, name)
	fmt.Printf("%s init\n", name)
	defer fmt.Printf("%s exit\n", name)
	handler(ctx)
}

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(16, WorkerInt(namedWorker))
	defer pool.Close()
	time.Sleep(time.Millisecond * 100)
	for i := 0; i < 100; i++ {
		taskId := i
		task := func(ctx context.Context) {
			name, ok := fromWorkerNameContext(ctx)
			if !ok {
				return
			}
			fmt.Printf("%s run: task %d\n", name, taskId)
			time.Sleep(time.Millisecond * 10)
		}
		if err := pool.RunContext(context.Background(), task); err != nil {
			panic(err)
		}
	}
}
