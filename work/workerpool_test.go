package work

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

var workId int32 = 0

type workNameKey struct{}

func newWorkNameContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, workNameKey{}, name)
}

func fromWorkNameContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(workNameKey{}).(string)
	return name, ok
}

func namedWorkProxy(ctx context.Context, handler WorkHandler) {
	name := fmt.Sprintf("work%d", atomic.AddInt32(&workId, 1))
	fmt.Printf("%s init\n", name)
	defer fmt.Printf("%s done\n", name)
	handler(newWorkNameContext(ctx, name))
}

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(16, WorkProxy(namedWorkProxy))
	defer pool.Close()
	time.Sleep(time.Millisecond * 100)
	for i := 0; i < 100; i++ {
		taskId := i
		task := func(ctx context.Context) {
			name, ok := fromWorkNameContext(ctx)
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
