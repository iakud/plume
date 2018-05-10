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
	done        chan struct{}
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
	this.done = make(chan struct{})
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
	defer close(this.done)
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
		connection := this.newConnection(conn)
		this.addConnection(connection)
		go this.serveConnection(connection)
	}
}

func (this *TCPServer) newConnection(conn *net.TCPConn) *TCPConnection {
	connection := newTCPConnection(conn)
	connection.connectFunc = this.ConnectFunc
	connection.disconnectFunc = this.DisconnectFunc
	connection.receiveFunc = this.ReceiveFunc
	return connection
}

func (this *TCPServer) serveConnection(connection *TCPConnection) {
	connection.serve()
	this.removeConnection(connection)
}

func (this *TCPServer) addConnection(connection *TCPConnection) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.connections[connection] = struct{}{}
}

func (this *TCPServer) removeConnection(connection *TCPConnection) {
	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.connections, connection)
}

func (this *TCPServer) closeConnections() {
	this.mu.Lock()
	defer this.mu.Unlock()
	for connection := range this.connections {
		connection.close()
		delete(this.connections, connection)
	}
}

func (this *TCPServer) Close() error {
	if this.listener != nil {
		return nil
	}
	err := this.listener.Close()
	this.listener = nil
	this.closeConnections()
	<-this.done
	return err
}
