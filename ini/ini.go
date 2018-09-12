package ini

import (
	"bufio"
	"io"
	"os"
)

func Open(name string) (*File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return parse(file)
}

func parse(r io.Reader) (*File, error) {
	iniFile := &File{}
	rd := bufio.NewReader(r)
	for {
		_, err := rd.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
	}

	return iniFile, nil
}
