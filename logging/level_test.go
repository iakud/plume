package logging

import (
	"fmt"
	"testing"
)

func TestLevel(t *testing.T) {
	fmt.Printf("Level num: %d\n", numLevel)
	fmt.Printf("%s: %d\n", TraceLevel, TraceLevel)
	fmt.Printf("%s: %d\n", DebugLevel, DebugLevel)
	fmt.Printf("%s: %d\n", InfoLevel, InfoLevel)
	fmt.Printf("%s: %d\n", WarningLevel, WarningLevel)
	fmt.Printf("%s: %d\n", ErrorLevel, ErrorLevel)
	fmt.Printf("%s: %d\n", PanicLevel, PanicLevel)
	fmt.Printf("%s: %d\n", FatalLevel, FatalLevel)
	fmt.Printf("%s Enabled(%s): %v\n", InfoLevel, DebugLevel, InfoLevel.Enabled(DebugLevel))
	fmt.Printf("%s Enabled(%s): %v\n", InfoLevel, InfoLevel, InfoLevel.Enabled(InfoLevel))
	fmt.Printf("%s Enabled(%s): %v\n", InfoLevel, WarningLevel, InfoLevel.Enabled(WarningLevel))
}
