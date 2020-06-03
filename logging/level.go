package logging

type Level int32

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
	numLevel
)

var levelName = [numLevel]string{
	TraceLevel: "TRACE ",
	DebugLevel: "DEBUG ",
	InfoLevel:  "INFO  ",
	WarnLevel:  "WARN  ",
	ErrorLevel: "ERROR ",
	PanicLevel: "PANIC ",
	FatalLevel: "FATAL ",
}

func (l Level) String() string {
	return levelName[l]
}

func (l Level) Enabled(level Level) bool {
	return level >= l
}

func (l Level) Disabled(level Level) bool {
	return level < l
}
