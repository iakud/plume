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
	sev    severity
	file   *os.File
	writer *bufio.Writer
}

func newLogWriter(logger *logging, sev severity) *logWriter {
	newWriter := new(logWriter)
	newWriter.writer = bufio.NewWriterSize(w, bufferSize)
	newWriter.done = make(chan struct{})
	go newWriter.periodicalFlush()
	return newWriter
}

func (this *logWriter) Write(p []byte) (n int, err error) {
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

func (this *logWriter) rotateFile(now time.Time) error {
	if this.file != nil {
		this.Flush()
		this.file.Close()
	}
	severityName[this.sev]

	this.writer = bufio.NewWriterSize(this.file, size)
}
