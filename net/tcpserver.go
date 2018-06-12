package net

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/iakud/falcon"
)

type TCPServer struct {
	loop  *falcon.EventLoop
	addr  string
	codec Codec

	ConnectFunc    func(*TCPConnection)
	DisconnectFunc func(*TCPConnection)
	ReceiveFunc    func(*TCPConnection, []byte)

	mutex       sync.Mutex
	listener    *net.TCPListener
	connections map[*TCPConnection]struct{}
	started     bool
	closed      bool
}

func NewTCPServer(loop *falcon.EventLoop, addr string, codec Codec) *TCPServer {
	server := &TCPServer{
		loop:  loop,
		addr:  addr,
		codec: codec,
	}
	return server
}

func (this *TCPServer) Start() {
	this.mutex.Lock()
	if this.started || this.closed {
		this.mutex.Unlock()
		return
	}
	this.started = true
	this.mutex.Unlock()

	go this.listenAndServe()
}

func (this *TCPServer) listenAndServe() {
	ln, err := listenTCP(this.addr)
	if err != nil {
		log.Printf("TCPServer: error: %v", err)
		return
	}
	this.serve(ln)
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

func (this *TCPServer) serve(ln *net.TCPListener) {
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
		connection := newTCPConnection(this.loop, conn, this.codec)
		connection.connectFunc = this.ConnectFunc
		connection.disconnectFunc = this.DisconnectFunc
		connection.receiveFunc = this.ReceiveFunc

		if !this.newConnection(connection) {
			connection.close()
			return
		}

		go this.serveConnection(connection)
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

func (this *TCPServer) serveConnection(connection *TCPConnection) {
	connection.serve()

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
