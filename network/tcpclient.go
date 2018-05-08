package network

import (
	"log"
	"net"
	"sync"
	"time"
)

type TCPClient struct {
	addr string

	closed     bool
	mu         sync.Mutex
	connection *TCPConnection
	done       chan struct{}
}

func NewTCPClient(addr string) *TCPClient {
	client := &TCPClient{
		addr: addr,
		done: make(chan struct{}),
	}
	return client
}

func (this *TCPClient) Start(handler TCPHandler) error {
	go this.serve(handler)
	return nil
}

func (this *TCPClient) serve(handler TCPHandler) {
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		if this.closed {
			return
		}
		conn, err := this.dial()
		if err != nil {
			if this.closed {
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
		connection := newTCPConnection(conn, handler)

		this.mu.Lock()
		if this.closed {
			this.mu.Unlock()
			conn.Close()
			return
		}
		this.connection = connection
		this.mu.Unlock()
		// this.setConnection(connection)

		this.serveConnection(connection)
	}
}

func (this *TCPClient) dial() (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", this.addr)
	if err != nil {
		return nil, err
	}
	return net.DialTCP("tcp", nil, raddr)

}

func (this *TCPClient) serveConnection(connection *TCPConnection) {
	connection.serve()
	this.setConnection(nil)
}

func (this *TCPClient) setConnection(connection *TCPConnection) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.connection = connection
}

func (this *TCPClient) GetConnection() *TCPConnection {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.connection
}

func (this *TCPClient) closeConnection() {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.connection != nil {
		this.connection.Close()
		this.connection = nil
	}
}

func (this *TCPClient) Close() error {
	this.closed = true
	this.closeConnection()
	<-this.done
	return nil
}
