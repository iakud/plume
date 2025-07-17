package network

import (
	"errors"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	ErrWSServerClosed = errors.New("network: WebSocket server closed")
)

type WSServer struct {
	Handler WSHandler

	mutex       sync.Mutex
	connections map[*WSConnection]struct{}
	closed      bool
}

func NewWSServer(handler WSHandler) *WSServer {
	server := &WSServer{
		Handler:     handler,
		connections: make(map[*WSConnection]struct{}),
	}
	return server
}

func WebsocketHandler(ws *WSServer) http.Handler {
	return websocket.Handler(ws.ServeWS)
}

func (s *WSServer) ServeWS(conn *websocket.Conn) {
	handler := s.Handler
	if handler == nil {
		handler = DefaultWSHandler
	}

	connection := newWSConnection(conn)
	if err := s.newConnection(connection); err != nil {
		connection.Close() // close
		return
	}
	s.serveConnection(connection, handler)
}

func (s *WSServer) newConnection(connection *WSConnection) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.closed {
		return ErrWSServerClosed
	}
	s.connections[connection] = struct{}{}
	return nil
}

func (s *WSServer) serveConnection(connection *WSConnection, handler WSHandler) {
	connection.serve(handler)
	// remove connection
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.closed {
		return
	}
	delete(s.connections, connection)
}

func (s *WSServer) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return
	}
	s.closed = true
	for connection := range s.connections {
		connection.Close()
		delete(s.connections, connection)
	}
}
