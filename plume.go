package plume

import (
	"net/http"
	_ "net/http/pprof"
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

func Run(app App) {
	if !atomic.CompareAndSwapInt32(&running, 0, 1) {
		// FIXME: log
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		sig := <-c
		_ = sig
		Close()
	}()
	go http.ListenAndServe(":8080", nil)

	<-done
	atomic.StoreInt32(&running, 0)
}

func Close() {
	close(done)
}
