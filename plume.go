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

type Service interface {
	Init()
	Shutdown()
}

func Run(s Service) {
	if !atomic.CompareAndSwapInt32(&running, 0, 1) {
		log.Info("Plume has running")
		return
	}
	log.Infof("Plume starting up")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go http.ListenAndServe(":80", nil)
	
	s.Init()
	select {
	case sig := <-c:
		log.Info("Plume got signal", sig)
	case <-ctx.Done():
	}
	log.Infof("Plume closing down")
	s.Shutdown()

	atomic.StoreInt32(&running, 0)
}

func Shutdown() {
	cancel()
}
