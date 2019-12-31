package work

import (
	"fmt"
	"testing"
	"time"
)

func wait() {
	time.Sleep(time.Second)
	fmt.Printf("wait()\n")
}

func TestWork(t *testing.T) {
	work := NewWork(wait)
	work.Join()
	fmt.Printf("work.Join()\n")
}
