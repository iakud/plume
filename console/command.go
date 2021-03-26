package console

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func Stack(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 1<<20)
	buf = buf[:runtime.Stack(buf, true)]
	filename := fileName("") + ".stack"
	if err := ioutil.WriteFile(filename, buf, 0666); err != nil {
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.WriteString(w, filename)
}

func fileName(name string) string {
	now := time.Now()
	return fmt.Sprintf("%s.%04d%02d%02d-%02d%02d%02d",
		name,
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())
}
