package net

import (
	"log"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	addr string

	mutex       sync.Mutex
	listener    *net.TCPListener
	connections map[*TCPConnection]struct{}
	started     bool
	closed      bool
}

func NewTCPServer(addr string) *TCPServer {
	server := &TCPServer{
		addr: addr,
	}
	return server
}

func (this *TCPServer) Start(connectionFunc func(*TCPConnection)) {
	this.mutex.Lock()
	if this.started || this.closed {
		this.mutex.Unlock()
		return
	}
	this.started = true
	this.mutex.Unlock()

	go this.listenAndServe(connectionFunc)
}

func (this *TCPServer) listenAndServe(connectionFunc func(*TCPConnection)) {
	ln, err := listenTCP(this.addr)
	if err != nil {
		log.Printf("TCPServer: error: %v", err)
		return
	}

	this.serve(ln, connectionFunc)
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

func (this *TCPServer) serve(ln *net.TCPListener, connectionFunc func(*TCPConnection)) {
	defer ln.Close()
	if !this.newListener(ln) {
		return
	}
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			if this.isClosed() {
				return
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
			return
		}
		tempDelay = 0
		connection := newTCPConnection(conn)

		if !this.newConnection(connection) {
			connection.close()
			return
		}

		go this.serveConnection(connection, connectionFunc)
	}
}

func (this *TCPServer) isClosed() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.closed
}

func (this *TCPServer) newListener(ln *net.TCPListener) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return false
	}
	this.listener = ln
	this.connections = make(map[*TCPConnection]struct{})
	return true
}

func (this *TCPServer) newConnection(connection *TCPConnection) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return false
	}
	this.connections[connection] = struct{}{}
	return true
}

func (this *TCPServer) removeConnection(connection *TCPConnection) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return
	}
	delete(this.connections, connection)
}

func (this *TCPServer) serveConnection(connection *TCPConnection, connectionFunc func(*TCPConnection)) {
	connectionFunc(connection)

	this.removeConnection(connection)
}

func (this *TCPServer) Close() {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.closed = true
	this.mutex.Unlock()

	if this.listener == nil {
		return
	}
	this.listener.Close()
	for connection := range this.connections {
		connection.close()
		delete(this.connections, connection)
	}
}
