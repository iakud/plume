package util

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var workerId int32 = 0

func workerInit(worker *Worker) {
	id := atomic.AddInt32(&workerId, 1)
	name := fmt.Sprintf("worker%d", id)
	worker.Context = name
	fmt.Printf("init: %s\n", name)
}

func TestWorkerPool(t *testing.T) {
	workerPool := NewWorkerPool(3, WorkerPoolSizeInfinite, workerInit)
	for i := 0; i < 100; i++ {
		buf := fmt.Sprintf("task %d", i)
		workerPool.Run(func(worker *Worker) {
			name := worker.Context.(string)
			fmt.Printf("%s %s\n", name, buf)
			time.Sleep(time.Millisecond * 100)
		})
	}
	var wg sync.WaitGroup
	wg.Add(1)
	workerPool.Run(func(*Worker) {
		wg.Done()
	})
	wg.Wait()
	workerPool.Close()
}
