package util

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	workerPool := NewWorkerPool(0, func(worker *Worker) {
		name := fmt.Sprintf("worker%d", rand.Intn(100))
		fmt.Printf("%s init\n", name)
		worker.Context = name
	})
	workerPool.Start(3)
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
	workerPool.Stop()
}
