package plume

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/iakud/plume/log"
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
		log.Info("Plume has running")
		return
	}
	log.Infof("Plume starting up")
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
	log.Infof("Plume closing down")
	app.Destory()
	atomic.StoreInt32(&running, 0)
}

func Close() {
	cancel()
}
