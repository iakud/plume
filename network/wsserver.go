package network

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrServerClosed = errors.New("network: WebSocket server closed")
)

type WSServer struct {
	mutex  sync.Mutex
	conns  map[*WSConn]struct{}
	closed bool
}

var upgrader = websocket.Upgrader{}

func NewWSServer() *WSServer {
	server := &WSServer{}
	server.conns = make(map[*WSConn]struct{})
	return server
}

func (s *WSServer) ServeWS(handler WSHandler, w http.ResponseWriter, r *http.Request) {
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
	s.serveConn(conn, handler)
}

func (s *WSServer) newConn(conn *WSConn) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.closed {
		return ErrServerClosed
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

func (s *WSServer) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return
	}
	s.closed = true
	for conn := range s.conns {
		conn.Close()
		delete(s.conns, conn)
	}
}
