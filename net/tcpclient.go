package net

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/iakud/falcon"
)

type TCPClient struct {
	loop *falcon.EventLoop
	addr string

	ConnectFunc    func(*TCPConnection)
	DisconnectFunc func(*TCPConnection)
	ReceiveFunc    func(*TCPConnection, []byte)

	mutex      sync.Mutex
	connection *TCPConnection
	started    bool
	closed     bool
}

func NewTCPClient(loop *falcon.EventLoop, addr string) *TCPClient {
	client := &TCPClient{
		loop: loop,
		addr: addr,
	}
	return client
}

func (this *TCPClient) Start() {
	this.mutex.Lock()
	if this.started || this.closed {
		this.mutex.Unlock()
		return
	}
	this.started = true
	this.mutex.Unlock()

	go this.connect()
}

func dialTCP(addr string) (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return net.DialTCP("tcp", nil, raddr)
}

func (this *TCPClient) connect() {
	var tempDelay time.Duration // how long to sleep on connect failure
	for {
		conn, err := dialTCP(this.addr)
		if err != nil {
			if this.isClosed() {
				return
			}

			if tempDelay == 0 {
				tempDelay = 1 * time.Second
			} else {
				tempDelay *= 2
			}
			if max := 1 * time.Minute; tempDelay > max {
				tempDelay = max
			}
			log.Printf("TCPClient: connect error: %v; retrying in %v", err, tempDelay)
			time.Sleep(tempDelay)
			continue
		}
		tempDelay = 0

		connection := newTCPConnection(this.loop, conn)

		if !this.newConnection(connection) {
			connection.close()
			return
		}

		go this.serveConnection(connection)
		return
	}
}

func (this *TCPClient) isClosed() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.closed
}

func (this *TCPClient) newConnection(connection *TCPConnection) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return false
	}
	this.connection = connection
	return true
}

func (this *TCPClient) removeConnection(connection *TCPConnection) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return false
	}
	this.connection = nil
	return true
}

func (this *TCPClient) serveConnection(connection *TCPConnection) {
	connection.serve(this.ConnectFunc, this.DisconnectFunc, this.ReceiveFunc)

	if this.removeConnection(connection) {
		go this.connect()
	}
}

func (this *TCPClient) GetConnection() *TCPConnection {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return nil
	}
	return this.connection
}

func (this *TCPClient) Close() {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.closed = true
	this.mutex.Unlock()

	if this.connection == nil {
		return
	}
	this.connection.close()
	this.connection = nil
}
