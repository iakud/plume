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

func (hooks Hooks) Add(hook Hook) {
	hooks = append(hooks, hook)
}

func (hooks Hooks) log(entry *Entry) {
	for _, hook := range hooks {
		if err := hook(entry); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to hook: %v\n", err)
		}
	}
}
