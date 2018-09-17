package log

import (
	"sync"
)

type AsyncLogger struct {
	messages []*message
	cond     *sync.Cond
}
