package plume

import (
	"os"
	"os/signal"
	"sync/atomic"
)

var running int32
var done = make(chan struct{})

type App interface {
	Init()
	Destory()
}

func Run() {
	if !atomic.CompareAndSwapInt32(&running, 0, 1) {
		// FIXME: log
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	<-done
	atomic.StoreInt32(&running, 0)
}

func Close() {
	close(done)
}
