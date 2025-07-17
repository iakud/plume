package network

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrWSServerClosed = errors.New("network: WebSocket server closed")
)

type WSServer struct {
	addr string
	handler WSHandler

	mutex  sync.Mutex
	server *http.Server
	conns  map[*WSConn]struct{}
	closed bool
}

var upgrader = websocket.Upgrader{}

func NewWSServer(addr string, handler WSHandler) *WSServer {
	if handler == nil {
		handler = DefaultWSHandler
	}
	server := &WSServer{
		addr: addr,
		handler: handler,
		conns: make(map[*WSConn]struct{}),
	}
	return server
}

func ListenAndServeWS(addr string, handler WSHandler) error {
	server := NewWSServer(addr, handler)
	return server.ListenAndServe()
}

func (s *WSServer) ListenAndServe() error {
	if s.isClosed() {
		return ErrServerClosed
	}

	s.server = &http.Server{Addr: s.addr, Handler: s}
	return s.server.ListenAndServe()
}

func (s *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn := newWSConn(wsconn)
	if err := s.newConn(conn); err != nil {
		conn.Close() // close
		return
	}
	s.serveConn(conn, s.handler)
}

func (s *WSServer) isClosed() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.closed
}

func (s *WSServer) newConn(conn *WSConn) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.closed {
		return ErrWSServerClosed
	}
	s.conns[conn] = struct{}{}
	return nil
}

func (s *WSServer) serveConn(conn *WSConn, handler WSHandler) {
	conn.serve(handler)
	// remove connection
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.closed {
		return
	}
	delete(s.conns, conn)
}

func (s *WSServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *WSServer) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return
	}
	s.closed = true
	if s.server == nil {
		return
	}
	s.server.Close()
	for conn := range s.conns {
		conn.Close()
		delete(s.conns, conn)
	}
}
