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
	numLevel = FatalLevel
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

func (this Level) String() string {
	return levelName[this]
}
