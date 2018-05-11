package network

import (
	"log"
	"net"
	"sync"
	"time"
)

type TCPClient struct {
	addr  string
	retry bool

	mu         sync.Mutex
	connection *TCPConnection
	closed     bool

	ConnectFunc    func(*TCPConnection)
	DisconnectFunc func(*TCPConnection)
	ReceiveFunc    func(*TCPConnection, []byte)
}

func NewTCPClient(addr string) *TCPClient {
	client := &TCPClient{
		addr:  addr,
		retry: true,
	}
	return client
}

func (this *TCPClient) Start() error {
	go this.serve()
	return nil
}

func (this *TCPClient) serve() {
	var tempDelay time.Duration // how long to sleep on connect failure
	for {
		conn, err := this.connect()
		if err != nil {
			this.mu.Lock()
			if this.closed {
				this.mu.Unlock()
				return
			}
			this.mu.Unlock()
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
		this.newConnection(conn)
	}
}

func (this *TCPClient) connect() (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", this.addr)
	if err != nil {
		return nil, err
	}
	return net.DialTCP("tcp", nil, raddr)
}

func (this *TCPClient) newConnection(conn *net.TCPConn) {
	connection := newTCPConnection(conn)
	connection.connectFunc = this.ConnectFunc
	connection.disconnectFunc = this.DisconnectFunc
	connection.receiveFunc = this.ReceiveFunc

	this.mu.Lock()
	if this.closed {
		this.mu.Unlock()
		connection.close()
		return
	}
	this.connection = connection
	this.mu.Unlock()

	this.serveConnection(connection)
}

func (this *TCPClient) serveConnection(connection *TCPConnection) {
	connection.serve()

	this.mu.Lock()
	defer this.mu.Unlock()
	if this.closed {
		return
	}
	this.connection = nil
}

func (this *TCPClient) GetConnection() *TCPConnection {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.connection
}

func (this *TCPClient) Close() error {
	this.mu.Lock()
	if this.closed {
		this.mu.Unlock()
		return nil
	}
	this.closed = true
	if connection := this.connection; connection != nil {
		this.connection = nil
		this.mu.Unlock()
		connection.close()
		return nil
	}
	this.mu.Unlock()
	return nil
}
