package logging

import (
	"fmt"
	"os"
)

type Hook func(Level, []byte) error

type hooks []Hook

func (h hooks) hook(l Level, p []byte) {
	for _, hook := range h {
		if err := hook(l, p); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to hook: %v\n", err)
		}
	}
}
