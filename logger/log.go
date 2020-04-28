package logger

type severity int32

const (
	debugLog severity = iota
	infoLog
	warningLog
	errorLog
	fatalLog
	numSeverity = 5
)

const severityChar = "DIWEF"

var severityName = []string{
	debugLog:   "DEBUG",
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

func (s *severity) String() string {
	return strconv.FormatInt(int64(*s), 10)
}

func Debug(args ...interface{}) {

}

func Debugf(format string, args ...interface{}) {

}

func Info(args ...interface{}) {

}

func Infof(format string, args ...interface{}) {

}

func Warning(args ...interface{}) {

}

func Warningf(format string, args ...interface{}) {

}

func Error(args ...interface{}) {

}

func Errorf(format string, args ...interface{}) {

}

func Fatal(args ...interface{}) {

}

func Fatalf(format string, args ...interface{}) {

}
