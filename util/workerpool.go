package util

import (
	"context"
	"log"
	"runtime"
	"sync"
)

type Worker struct {
	Context interface{}

	ctx context.Context
}

func newWorker(workerFunc func(*Worker)) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	worker := &Worker{ctx: ctx}
	go func() {
		defer cancel()
		workerFunc(worker)
	}()
	return worker
}

func (this *Worker) join() {
	<-this.ctx.Done()
}

type WorkerPool struct {
	maxSize  int
	initFunc func(*Worker)

	mutex    sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
	running  bool
	workers  []*Worker
	queue    []func(*Worker)
}

func NewWorkerPool(maxSize int, initFunc func(*Worker)) *WorkerPool {
	pool := &WorkerPool{
		maxSize:  maxSize,
		initFunc: initFunc,
		running:  false,
	}
	pool.notFull = sync.NewCond(&pool.mutex)
	pool.notEmpty = sync.NewCond(&pool.mutex)
	return pool
}

func (this *WorkerPool) Start(numWorkers int) {
	this.mutex.Lock()
	if this.running {
		return
	}
	this.running = true
	var workers []*Worker
	for i := 0; i < numWorkers; i++ {
		worker := newWorker(this.workerRoutine)
		workers = append(workers, worker)
	}
	this.workers = workers
	this.mutex.Unlock()
}

func (this *WorkerPool) Stop() {
	this.mutex.Lock()
	this.running = false
	this.notEmpty.Broadcast()
	workers := this.workers
	this.workers = nil
	this.mutex.Unlock()
	for _, worker := range workers {
		worker.join()
	}
}

func (this *WorkerPool) Run(task func(*Worker)) {
	this.mutex.Lock()
	for this.maxSize > 0 && len(this.queue) >= this.maxSize {
		this.notFull.Wait()
	}
	this.queue = append(this.queue, task)
	this.mutex.Unlock()
	this.notEmpty.Signal()
}

func (this *WorkerPool) workerRoutine(worker *Worker) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("util: panic worker: %v\n%s", err, buf)
		}
	}()
	if initFunc := this.initFunc; initFunc != nil {
		initFunc(worker)
	}
	var running bool = true
	for running {
		this.mutex.Lock()
		for this.running && len(this.queue) == 0 {
			this.notEmpty.Wait()
		}
		var task func(*Worker)
		var notFull bool
		if len(this.queue) > 0 {
			task = this.queue[0]
			this.queue = this.queue[1:]
			notFull = this.maxSize > 0
		}
		running = this.running
		this.mutex.Unlock()
		if notFull {
			this.notFull.Signal()
		}
		if task != nil {
			task(worker)
		}
	}
}
