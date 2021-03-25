package plume

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	c := make(chan os.Signal, 1)
	// signal.Notify(c, signal.)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer stop()
	<-ctx.Done()
}
