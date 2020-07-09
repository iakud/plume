package logging

import (
	"fmt"
	"os"
	"time"
)

type Entry struct {
	Time    time.Time
	Level   Level
	Message string
	PC      uintptr
	File    string
	Line    int
}

type Hook func(*Entry) error

type Hooks []Hook

func (hooks Hooks) hook(e *Entry) {
	for _, hook := range hooks {
		if err := hook(e); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to hook: %v\n", err)
		}
	}
}
