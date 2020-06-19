package logging

import (
	"fmt"
)

type Level int32

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarningLevel
	ErrorLevel
	PanicLevel
	FatalLevel
	numLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

func (l Level) Enabled(level Level) bool {
	return level >= l
}

func (l Level) Disabled(level Level) bool {
	return level < l
}
