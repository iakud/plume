package logging

import (
	"os"
	"testing"
)

func hookWarning(l Level, b []byte) error {
	if WarningLevel.Enabled(l) {
		os.Stdout.Write(b)
	}
	return nil
}

func TestHook(t *testing.T) {
	hooks := hooks{hookWarning}
	hooks.hook(InfoLevel, []byte("Info log!"))
	hooks.hook(WarningLevel, []byte("Warning log!"))
	hooks.hook(ErrorLevel, []byte("Error log!"))
}
