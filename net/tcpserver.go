package net

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var (
	ErrServerClosed = errors.New("net: Server closed")
)

type TCPServer struct {
	addr    string
	handler TCPHandler
	codec   Codec

	mutex       sync.Mutex
	listener    *net.TCPListener
	connections map[*TCPConnection]struct{}
	closed      bool
}

func NewTCPServer(addr string, handler TCPHandler, codec Codec) *TCPServer {
	server := &TCPServer{
		addr:    addr,
		handler: handler,
		codec:   codec,
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

func (this *TCPServer) ListenAndServe() error {
	if this.isClosed() {
		return ErrServerClosed
	}
	ln, err := listenTCP(this.addr)
	if err != nil {
		return err
	}
	return this.serve(ln)
}

func (this *TCPServer) serve(ln *net.TCPListener) error {
	defer ln.Close()

	if err := this.newListener(ln); err != nil {
		return err
	}
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			if this.isClosed() {
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
				log.Printf("TCPServer: accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			log.Printf("TCPServer: error: %v", err)
			return err
		}
		tempDelay = 0

		connection := newTCPConnection(conn, this.handler, this.codec)
		if err := this.newConnection(connection); err != nil {
			connection.close() // close
			return err
		}
		go this.serveConnection(connection)
	}
}

func (this *TCPServer) isClosed() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.closed
}

func (this *TCPServer) newListener(ln *net.TCPListener) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return ErrServerClosed
	}
	this.listener = ln
	this.connections = make(map[*TCPConnection]struct{})
	return nil
}

func (this *TCPServer) newConnection(connection *TCPConnection) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return ErrServerClosed
	}
	this.connections[connection] = struct{}{}
	return nil
}

func (this *TCPServer) serveConnection(connection *TCPConnection) {
	connection.serve()
	// remove connection
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return
	}
	delete(this.connections, connection)
}

func (this *TCPServer) Close() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return nil
	}
	this.closed = true
	if this.listener == nil {
		return nil
	}
	err := this.listener.Close()
	this.listener = nil
	for connection := range this.connections {
		connection.close()
		delete(this.connections, connection)
	}
	return err
}
