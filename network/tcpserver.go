package network

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var (
	ErrServerClosed = errors.New("network: Server closed")
)

type TCPServer struct {
	addr string

	mutex       sync.Mutex
	listener    *net.TCPListener
	connections map[*TCPConnection]struct{}
	closed      bool
}

func NewTCPServer(addr string) *TCPServer {
	server := &TCPServer{
		addr: addr,
	}
	return server
}

func listenTCP(addr string) (*net.TCPListener, error) {
	if addr == "" {
		addr = ":0"
	}
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", laddr)
}

func (s *TCPServer) ListenAndServe(handler TCPHandler, codec Codec) error {
	if s.isClosed() {
		return ErrServerClosed
	}
	ln, err := listenTCP(s.addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	if err := s.newListener(ln); err != nil {
		return err
	}

	if handler == nil {
		handler = DefaultTCPHandler
	}
	if codec == nil {
		codec = DefaultCodec
	}

	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			if s.isClosed() {
				return ErrServerClosed
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("network: TCPServer accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			log.Printf("network: TCPServer error: %v", err)
			return err
		}
		tempDelay = 0

		connection := newTCPConnection(conn)
		if err := s.newConnection(connection); err != nil {
			connection.Close() // close
			return err
		}
		go s.serveConnection(connection, handler, codec)
	}
}

func (s *TCPServer) isClosed() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.closed
}

func (s *TCPServer) newListener(ln *net.TCPListener) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return ErrServerClosed
	}
	s.listener = ln
	s.connections = make(map[*TCPConnection]struct{})
	return nil
}

func (s *TCPServer) newConnection(connection *TCPConnection) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return ErrServerClosed
	}
	s.connections[connection] = struct{}{}
	return nil
}

func (s *TCPServer) serveConnection(connection *TCPConnection, handler TCPHandler, codec Codec) {
	connection.serve(handler, codec)
	// remove connection
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return
	}
	delete(s.connections, connection)
}

func (s *TCPServer) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return
	}
	s.closed = true
	if s.listener == nil {
		return
	}
	s.listener.Close()
	s.listener = nil
	for connection := range s.connections {
		connection.Close()
		delete(s.connections, connection)
	}
}
