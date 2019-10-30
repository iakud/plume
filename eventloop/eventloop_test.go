package eventloop

import (
	"fmt"
	"testing"
	"time"
)

func TestEventLoop(t *testing.T) {
	loop := NewEventLoop()
	loop.RunInLoop(func() {
		fmt.Println("close in loop")
		loop.Close()
	})
	loop.Loop()
}

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

func TestTicker(t *testing.T) {
	loop := NewEventLoop()
	times := 0
	ticker := loop.RunEvery(time.Second, func() {
		times++
		fmt.Println("on ticker", times)
	})
	defer ticker.Stop()
	loop.RunAfter(time.Second*3, func() {
		ticker.Stop()
	})
	countDown := 5
	loop.RunEvery(time.Second, func() {
		countDown--
		fmt.Println("count down", countDown)
		if countDown == 0 {
			loop.Close()
			return
		}
	})
	loop.Loop()
}
