package util

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLockGuard(t *testing.T) {
	var l sync.Mutex
	go func() {
		LockGuard(&l, func() {
			fmt.Println("lock begin")
			time.Sleep(time.Second)
			fmt.Println("lock end")
		})
	}()
	time.Sleep(time.Millisecond)
	LockGuard(&l, func() {
		fmt.Println("in lock")
	})

	var rwl sync.RWMutex
	go func() {
		// rlock
		LockGuard(rwl.RLocker(), func() {
			fmt.Println("rlock begin")
			time.Sleep(time.Second)
			fmt.Println("rlock end")
		})
	}()

	time.Sleep(time.Millisecond)
	// rlock
	LockGuard(rwl.RLocker(), func() {
		fmt.Println("in rlock")
	})
	// lock
	LockGuard(&rwl, func() {
		fmt.Println("in wlock")
	})
}
