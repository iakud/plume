package plume

import (
	"os"
	"os/signal"
)

func Run() {
	c := make(chan os.Signal, 1)
	// signal.Notify(c, signal.)
	os.Interrupt
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer stop()
	<-ctx.Done()
}
