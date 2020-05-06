package logger

const bufferSize = 256 * 1024

type rollingFileWriter struct {
}

type rollingFileWriterTime struct {
	*rollingFileWriter
	name string
	file *os.File

	startOfPeriod time.Time
}

func newFileWriter(logger *logging, path string, name string) *fileWriter {
	fileWriter := new(fileWriter)
	fileWriter.writer = bufio.NewWriterSize(w, bufferSize)
	return newWriter
}

func (this *fileWriter) Write(p []byte) (n int, err error) {
	thisPeriod := time.Now().Truncate(time.Hour)
	if thisPeriod != this.startOfPeriod {
		this.rollFile(now)
	}
	return this.writer.Write(p)
}

func (this *fileWriter) Flush() error {
	return this.writer.Flush()
}

func (this *fileWriter) Sync() error {
	return this.file.Sync()
}

func (this *fileWriter) rollFile(now time.Time) error {
	if this.file != nil {
		this.Flush()
		this.file.Close()
	}
	this.startOfPeriod = time.Now().Truncate(time.Hour)

	this.writer = bufio.NewWriterSize(this.file, size)
}

func (this *fileWriter) createFile(filename string, now time.Time) error {
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	this.file = file
	this.writer = bufio.NewWriterSize(file, bufferSize)
	return nil
}
