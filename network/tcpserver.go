package network

import (
	"log"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	addr string

	ConnectFunc    func(*TCPConnection)
	DisconnectFunc func(*TCPConnection)
	ReceiveFunc    func(*TCPConnection, []byte)

	mu          sync.Mutex
	connections map[*TCPConnection]struct{}
	listener    *net.TCPListener
	closed      bool
}

func NewTCPServer(addr string) *TCPServer {
	server := &TCPServer{
		addr:        addr,
		connections: make(map[*TCPConnection]struct{}),
	}
	return server
}

func (this *TCPServer) Start() error {
	if this.listener != nil {
		return nil
	}
	ln, err := this.listen()
	if err != nil {
		return err
	}
	this.listener = ln
	go this.serve()
	return nil
}

func (this *TCPServer) listen() (*net.TCPListener, error) {
	addr := this.addr
	if addr == "" {
		addr = ":0"
	}
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", laddr)
}

func (this *TCPServer) serve() {
	defer this.listener.Close()

	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := this.listener.AcceptTCP()
		if err != nil {
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
		this.newConnection(conn)
	}
}

func (this *TCPServer) newConnection(conn *net.TCPConn) {
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
	this.connections[connection] = struct{}{}
	this.mu.Unlock()

	go this.serveConnection(connection)
}

func (this *TCPServer) serveConnection(connection *TCPConnection) {
	connection.serve()

	this.mu.Lock()
	defer this.mu.Unlock()
	if this.closed {
		return
	}
	delete(this.connections, connection)
}

func (this *TCPServer) Close() error {
	if this.listener != nil {
		return nil
	}

	this.mu.Lock()
	if this.closed {
		this.mu.Unlock()
		return nil
	}
	this.closed = true
	var err error
	if this.listener != nil {
		err = this.listener.Close()
	}
	for connection := range this.connections {
		connection.close()
		delete(this.connections, connection)
	}
	return err
}
