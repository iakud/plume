package logger

type Level int32

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	numLevel = 6
)

var levelName = [numLevel]string{
	TraceLevel: "TRACE ",
	DebugLevel: "DEBUG ",
	InfoLevel:  "INFO  ",
	WarnLevel:  "WARN  ",
	ErrorLevel: "ERROR ",
	FatalLevel: "FATAL ",
}

func (this Level) String() string {
	return levelName[this]
}
