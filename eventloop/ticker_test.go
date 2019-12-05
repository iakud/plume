package eventloop

import (
	"fmt"
	"testing"
	"time"
)

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
