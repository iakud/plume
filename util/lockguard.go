package util

import (
	"sync"
)

func LockGuard(locker sync.Locker, f func()) {
	locker.Lock()
	defer locker.Unlock()
	f()
}
