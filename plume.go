package plume

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync/atomic"
)

var running int32
var ctx, cancel = context.WithCancel(context.Background())

type App interface {
	Init()
	Run(context.Context)
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
	app.Init()
	app.Run(ctx)
	app.Destory()
	atomic.StoreInt32(&running, 0)
}

func Close() {
	cancel()
}
