package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func logName(name string, t time.Time) string {
	logname := fmt.Sprintf("%s.%04d%02d%02d-%02d%02d%02d",
		name,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second())
	return logname
}

func createFile(dir, name string, t time.Time) (*os.File, error) {
	logname := logName(name, t)
	filename := filepath.Join(dir, logname)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("logging: cannot create log: %v", err)
	}

	symlink := filepath.Join(dir, name)
	os.Remove(symlink) // ignore err
	if err := os.Symlink(logname, symlink); err != nil {
		os.Link(logname, symlink)
	}
	return file, nil
}
