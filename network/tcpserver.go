package network

import (
	"log"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	addr string

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

func (this *TCPServer) Start(handler TCPHandler) error {
	if this.listener != nil {
		return nil
	}
	addr := this.addr
	if addr == "" {
		addr = ":0"
	}
	la, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp", la)
	if err != nil {
		return err
	}
	this.listener = ln
	this.done = make(chan struct{})
	go this.serve(handler)
	return nil
}

func (this *TCPServer) serve(handler TCPHandler) {
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
		connection := newTCPConnection(conn, handler)
		this.addConnection(connection)
		go this.serveConnection(connection)
	}
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
		connection.Close()
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
