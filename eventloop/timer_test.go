package eventloop

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	loop := NewEventLoop()
	loop.RunAfter(time.Second, func() {
		fmt.Println("on timer1")
	})
	timer2 := loop.RunAfter(time.Second, func() {
		fmt.Println("on timer1, close")
		loop.Close()
	})
	timer2.Stop()
	loop.RunAfter(time.Second*2, func() {
		fmt.Println("on timer3, reset timer2")
		timer2.Reset(time.Second)
	})
	loop.Loop()
}
