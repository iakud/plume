package network

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	ErrWSServerClosed = errors.New("network: WebSocket server closed")
)

type WSServer struct {
	Handler WSHandler
	server  websocket.Server

	mutex       sync.Mutex
	connections map[*WSConnection]struct{}
	closed      bool
}

func NewWSServer(handler WSHandler) *WSServer {
	server := &WSServer{
		Handler:     handler,
		connections: make(map[*WSConnection]struct{}),
	}
	server.server = websocket.Server{Handler: server.serveWebSocket, Handshake: checkOrigin}
	return server
}

func checkOrigin(config *websocket.Config, req *http.Request) (err error) {
	config.Origin, err = websocket.Origin(config, req)
	if err == nil && config.Origin == nil {
		return fmt.Errorf("null origin")
	}
	return err
}

func (s *WSServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.server.ServeHTTP(w, req)
}

func (s *WSServer) serveWebSocket(conn *websocket.Conn) {
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
