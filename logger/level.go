package logger

type logLevel int32

const (
	traceLog logLevel = iota
	debugLog
	infoLog
	warningLog
	errorLog
	fatalLog
	numLevel = 6
)

var logLevelName = [numLevel]string{
	traceLog: "TRACE ",
	debugLog: "DEBUG ",
	infoLog:  "INFO  ",
	warnLog:  "WARN  ",
	errorLog: "ERROR ",
	fatalLog: "FATAL ",
}

func (this logLevel) String() string {
	return logLevelName[this]
}
