package logger

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
