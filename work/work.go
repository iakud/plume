package work

import (
	"context"
	"log"
	"runtime"
)

type Work struct {
	workFunc func()
	ctx      context.Context
}

func NewWork(workFunc func()) *Work {
	ctx, cancel := context.WithCancel(context.Background())
	work := &Work{
		workFunc: workFunc,
		ctx:      ctx,
	}
	go func() {
		defer cancel()
		work.runner()
	}()
	return work
}

func (this *Work) Join() {
	<-this.ctx.Done()
}

func (this *Work) runner() {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("work: panic work: %v\n%s", err, buf)
		}
	}()
	this.workFunc()
}
