package work

import (
	"log"
	"runtime"
)

type Worker struct {
	f    func()
	done chan struct{}
}

func NewWorker(f func()) *Worker {
	worker := &Worker{
		f:    f,
		done: make(chan struct{}),
	}
	go worker.runner()
	return worker
}

func (w *Worker) Done() <-chan struct{} {
	return w.done
}

func (w *Worker) runner() {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("work: panic runner: %v\n%s", err, buf)
		}
		close(w.done)
	}()
	w.f()
}
