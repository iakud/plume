package logging

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

const kCallerSkip int = 3

type Logger struct {
	level Level
	pool  *BufferPool
}

func New() *Logger {
	logger := &Logger{
		level: TraceLevel,
		pool:  NewBufferPool(),
	}
	return logger
}

func (logger *Logger) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&logger.level), int32(level))
}

func (logger *Logger) GetLevel() Level {
	return Level(atomic.LoadInt32((*int32)(&logger.level)))
}

func (logger *Logger) IsLevelDisabled(level Level) bool {
	return logger.GetLevel() > level
}

func (logger *Logger) output(l Level, b *bytes.Buffer) {
	if l > ErrorLevel {
		// FIXME: Sync()
	}
	if b.Bytes()[b.Len()-1] != '\n' {
		b.WriteByte('\n')
	}
	os.Stdout.Write(b.Bytes())

	logger.pool.Put(b)

	if FatalLevel == l {
		os.Exit(1)
	}
}

func fourDigits(buf *[]byte, i, d int) {
	(*buf)[i+3] = byte('0' + d%10)
	d /= 10
	(*buf)[i+2] = byte('0' + d%10)
	d /= 10
	(*buf)[i+1] = byte('0' + d%10)
	d /= 10
	(*buf)[i] = byte('0' + d%10)
}

func twoDigits(buf *[]byte, i, d int) {
	(*buf)[i+1] = byte('0' + d%10)
	d /= 10
	(*buf)[i] = byte('0' + d%10)
}

func someDigits(buf *[]byte, i, d int) int {
	var b [20]byte
	j := len(b) - 1
	for d >= 10 {
		b[j] = byte('0' + d%10)
		d /= 10
		j--
	}
	b[j] = byte('0' + d%10)
	return copy((*buf)[i:], b[j:])
}

func (logger *Logger) header(l Level) *bytes.Buffer {
	_, file, line, ok := runtime.Caller(kCallerSkip)
	if !ok {
		file = "???"
		line = 1
	} else {
		if slash := strings.LastIndexByte(file, '/'); slash >= 0 {
			file = file[slash+1:]
		}
	}
	return logger.formatHeader(l, file, line)
}

func (logger *Logger) formatHeader(l Level, file string, line int) *bytes.Buffer {
	now := time.Now()
	b := logger.pool.Get()
	buf := make([]byte, 32)
	// format yyyy/mm/dd hh:mm:ss level file:line:
	year, month, day := now.Date()
	fourDigits(&buf, 0, year)
	buf[4] = '/'
	twoDigits(&buf, 5, int(month))
	buf[7] = '/'
	twoDigits(&buf, 8, day)
	buf[10] = ' '
	hour, minute, second := now.Clock()
	twoDigits(&buf, 11, hour)
	buf[13] = ':'
	twoDigits(&buf, 14, minute)
	buf[16] = ':'
	twoDigits(&buf, 17, second)
	buf[19] = ' '
	b.Write(buf[:20])
	b.WriteString(l.String())
	b.WriteByte(' ')
	b.WriteString(file)
	buf[0] = ':'
	n := someDigits(&buf, 1, line)
	buf[n+1] = ':'
	buf[n+2] = ' '
	b.Write(buf[:n+3])
	return b
}

func (logger *Logger) logf(l Level, format string, a ...interface{}) {
	if logger.level.Enabled(l) {
		b := logger.header(l)
		fmt.Fprintf(b, format, a...)
		logger.output(l, b)
	}
}

func (logger *Logger) log(l Level, a ...interface{}) {
	if logger.level.Enabled(l) {
		b := logger.header(l)
		fmt.Fprint(b, a...)
		logger.output(l, b)
	}
}

func (logger *Logger) Tracef(format string, a ...interface{}) {
	logger.logf(TraceLevel, format, a...)
}

func (logger *Logger) Trace(a ...interface{}) {
	logger.log(TraceLevel, a...)
}

func (logger *Logger) Debugf(format string, a ...interface{}) {
	logger.logf(DebugLevel, format, a...)
}

func (logger *Logger) Debug(a ...interface{}) {
	logger.log(DebugLevel, a...)
}

func (logger *Logger) Infof(format string, a ...interface{}) {
	logger.logf(InfoLevel, format, a...)
}

func (logger *Logger) Info(a ...interface{}) {
	logger.log(InfoLevel, a...)
}

func (logger *Logger) Warningf(format string, a ...interface{}) {
	logger.logf(WarningLevel, format, a...)
}

func (logger *Logger) Warning(a ...interface{}) {
	logger.log(WarningLevel, a...)
}

func (logger *Logger) Errorf(format string, a ...interface{}) {
	logger.logf(ErrorLevel, format, a...)
}

func (logger *Logger) Error(a ...interface{}) {
	logger.log(ErrorLevel, a...)
}

func (logger *Logger) Panicf(format string, a ...interface{}) {
	logger.logf(PanicLevel, format, a...)
}

func (logger *Logger) Panic(a ...interface{}) {
	logger.log(PanicLevel, a...)
}

func (logger *Logger) Fatalf(format string, a ...interface{}) {
	logger.logf(FatalLevel, format, a...)
}

func (logger *Logger) Fatal(a ...interface{}) {
	logger.log(FatalLevel, a...)
}
