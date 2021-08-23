package plume

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/iakud/plume/log"
	"github.com/iakud/plume/service"
)

var running int32
var ctx, cancel = context.WithCancel(context.Background())

func Run(o ...Option) {
	if !atomic.CompareAndSwapInt32(&running, 0, 1) {
		log.Info("Plume has running")
		return
	}

	// options
	opts := options{}
	for _, option := range o {
		option(&opts)
	}

	log.Infof("Plume starting up")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go http.ListenAndServe(":80", nil)
	
	service.Init(opts.services)
	select {
	case sig := <-c:
		log.Info("Plume got signal", sig)
	case <-ctx.Done():
	}
	log.Infof("Plume closing down")
	
	service.Shutdown(opts.services)
	atomic.StoreInt32(&running, 0)
}

func Shutdown() {
	cancel()
}
