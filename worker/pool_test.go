package worker

import (
	"context"
	"fmt"
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

type namedPool struct {
}

func (this *namedPool) WorkerContext(ctx context.Context) context.Context {
	name := fmt.Sprintf("worker%d", atomic.AddInt32(&workerId, 1))
	fmt.Printf("%s init\n", name)
	return newWorkerNameContext(ctx, name)
}

func (this *namedPool) WorkerExit(ctx context.Context) {
	if name, ok := fromWorkerNameContext(ctx); ok {
		fmt.Printf("%s exit\n", name)
	}
}

func TestPool(t *testing.T) {
	pool := NewPool(3, 16, &namedPool{})
	defer pool.Close()
	time.Sleep(time.Second)
	for i := 0; i < 100; i++ {
		buf := fmt.Sprintf("task %d", i)
		task := func(ctx context.Context) {
			name, ok := fromWorkerNameContext(ctx)
			if !ok {
				return
			}
			fmt.Printf("%s run: %s\n", name, buf)
			time.Sleep(time.Millisecond * 100)
		}
		if err := pool.Run(context.Background(), task); err != nil {
			panic(err)
		}
	}
}
