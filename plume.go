package plume

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"

	"github.com/iakud/plume/log"
)

var running int32
var ctx, cancel = context.WithCancel(context.Background())

type Service interface {
	Init()
	Run(context.Context)
	Destory()
}

func Run(services ...Service) {
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
	/*
		s.Init()
		s.Run(ctx)
		log.Infof("Plume closing down")
		s.Destory()
	*/
	var stops []func()
	for _, s := range services {
		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		s.Init()
		go func() {
			s.Run(ctx)
			wg.Done()
		}()
		stop := func() {
			cancel()
			wg.Wait()
			s.Destory()
		}
		stops = append(stops, stop)
	}
	<-ctx.Done()
	log.Infof("Plume closing down")
	// destory
	for i := len(stops) - 1; i >= 0; i-- {
		stops[i]()
	}

	atomic.StoreInt32(&running, 0)
}

func Close() {
	cancel()
}
