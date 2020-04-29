package log

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

// bufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const bufferSize = 256 * 1024

type logWriter struct {
	logger *logging
	level  logLevel
	file   *os.File
	writer *bufio.Writer

	startOfPeriod time.Time
}

func newLogWriter(logger *logging, level logLevel) *logWriter {
	newWriter := new(logWriter)
	newWriter.writer = bufio.NewWriterSize(w, bufferSize)
	newWriter.done = make(chan struct{})
	go newWriter.periodicalFlush()
	return newWriter
}

func (this *logWriter) Write(p []byte) (n int, err error) {
	thisPeriod := time.Now().Truncate(time.Hour)
	if thisPeriod != this.startOfPeriod {
		this.rollFile(now)
	}

	n, err = this.writer.Write(p)
	if err != nil {
		this.logger.exit(err)
	}
	return n, err
}

func (this *logWriter) Flush() error {
	return this.writer.Flush()
}

func (this *logWriter) Sync() error {
	return this.file.Sync()
}

func (this *logWriter) rollFile(now time.Time) error {
	if this.file != nil {
		this.Flush()
		this.file.Close()
	}

	this.startOfPeriod = time.Now().Truncate(time.Hour)
	severityName[this.sev]

	this.writer = bufio.NewWriterSize(this.file, size)
}

func (this *logWriter) createFile(filename string, now time.Time) {
	this.file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
}
