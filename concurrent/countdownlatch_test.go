package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func TestCountDownLatch(t *testing.T) {
	const kCount = 3
	latch := NewCountDownLatch(kCount)
	for i := 0; i < 3; i++ {
		go func(n int) {
			time.Sleep(time.Millisecond * 10)
			fmt.Println("countdown", n)
			latch.CountDown()
		}(i)
	}
	latch.Wait()
	fmt.Println("down")
}
