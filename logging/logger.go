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

const kCallerSkip int = 2

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

func (logger *Logger) output(l Level, buf []byte) {
	if l > ErrorLevel {
		// FIXME: Sync()
	}
	os.Stdout.Write(buf)

	switch l {
	case PanicLevel:
		panic("FIXME")
	case FatalLevel:
		os.Exit(1)
	}
}

const digits = "0123456789"

func (logger *Logger) formatHeader(b *bytes.Buffer, l Level, caller *Caller) {
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	var buf [64]byte
	buf[0] = digits[(year/1000)%10]
	buf[1] = digits[(year/100)%10]
	buf[2] = digits[(year/10)%10]
	buf[3] = digits[year%10]
	buf[4] = digits[(month/10)%10]
	buf[5] = digits[month%10]
	buf[6] = digits[(day/10)%10]
	buf[7] = digits[day%10]
	buf[8] = ' '
	buf[9] = digits[(hour/10)%10]
	buf[10] = digits[hour%10]
	buf[11] = ':'
	buf[12] = digits[(minute/10)%10]
	buf[13] = digits[minute%10]
	buf[14] = ':'
	buf[15] = digits[(second/10)%10]
	buf[16] = digits[second%10]
	buf[17] = ' '
	b.Write(buf[:18])
	b.WriteString(l.String())
	b.WriteByte(' ')
	b.WriteString(caller.File)
	line := caller.Line
	i := 61
	buf[63] = ' '
	buf[62] = ':'
	for i >= 0 {
		buf[i] = digits[line%10]
		i--
		line /= 10
		if line == 0 {
			break
		}
	}
	buf[i] = ':'
	b.Write(buf[i:])
}

func someDigits(buf []byte, i int, num int) int {
	var n int = 0
	if num >= 10 {
		n = someDigits(buf, i, num/10)
	}
	buf[i+n] = digits[n%10]
	return n + 1
}

func (logger *Logger) logf(l Level, format string, a ...interface{}) {
	if logger.level.Disabled(l) {
		return
	}
	caller := newCaller(runtime.Caller(kCallerSkip))
	buffer := logger.pool.Get()
	// write file
	logger.formatHeader(buffer, l, &caller)
	fmt.Fprintf(buffer, format, a...)
	logger.output(l, buffer.Bytes())
	logger.pool.Put(buffer)
}

func (logger *Logger) Tracef(format string, a ...interface{}) {
	logger.logf(TraceLevel, format, a...)
}

func (logger *Logger) Debugf(format string, a ...interface{}) {
	logger.logf(DebugLevel, format, a...)
}

func (logger *Logger) Infof(format string, a ...interface{}) {
	logger.logf(InfoLevel, format, a...)
}

func (logger *Logger) Warnf(format string, a ...interface{}) {
	logger.logf(WarnLevel, format, a...)
}

func (logger *Logger) Errorf(format string, a ...interface{}) {
	logger.logf(ErrorLevel, format, a...)
}

func (logger *Logger) Panicf(format string, a ...interface{}) {
	logger.logf(PanicLevel, format, a...)
}

func (logger *Logger) Fatalf(format string, a ...interface{}) {
	logger.logf(FatalLevel, format, a...)
}

type Caller struct {
	PC   uintptr
	File string
	Line int
}

func newCaller(pc uintptr, file string, line int, ok bool) Caller {
	if !ok {
		file = "???"
		line = 1
	} else {
		if slash := strings.LastIndexByte(file, '/'); slash >= 0 {
			file = file[slash+1:]
		}
	}
	return Caller{
		PC:   pc,
		File: file,
		Line: line,
	}
}
