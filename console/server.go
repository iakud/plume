package console

import (
	"net/http"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sync"
)

var mux = http.NewServeMux()

type Server struct {
	server *http.Server
	mux    *http.ServeMux
}

func NewServer(addr string) *Server {
	server := &http.Server{Addr: addr}
	s := &Server{server: server, mux: http.NewServeMux()}
	return s
}

func (s *Server) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {

		}
	}()
	runtime.Stack()
}

func (s *Server) Stop() {
	s.server.Close()
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	mux.Handle(pattern, handler)
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.HandleFunc(pattern, handler)
}
