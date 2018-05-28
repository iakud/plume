package network

import (
	"bufio"
	"net"
)

type TCPConnection struct {
	conn  *net.TCPConn
	codec Codec

	sendQueue *SendQueue

	connectFunc    func(*TCPConnection)
	disconnectFunc func(*TCPConnection)
	receiveFunc    func(*TCPConnection, []byte)
}

func newTCPConnection(conn *net.TCPConn) *TCPConnection {
	connection := &TCPConnection{
		conn:      conn,
		codec:     &StdCodec{},
		sendQueue: NewSendQueue(),
	}
	conn.SetNoDelay(true) // no delay
	return connection
}

func (this *TCPConnection) serve() {
	defer this.conn.Close()

	this.onConnect()         // on connect
	defer this.onDisonnect() // on disconnect

	done := make(chan struct{})
	go this.serveWrite(done) // write

	wait := func() { <-done }
	defer wait() // wait write

	this.serveRead() // read
}

func (this *TCPConnection) serveRead() {
	defer this.conn.Close()
	defer this.sendQueue.Close()

	rd := bufio.NewReader(this.conn)
	for {
		b, err := this.codec.Read(rd)
		if err != nil {
			return
		}
		this.onReceive(b)
	}
}

func (this *TCPConnection) serveWrite(done chan struct{}) {
	defer close(done)
	defer this.conn.Close()

	w := bufio.NewWriter(this.conn)
	for {
		buffers := this.sendQueue.Get()
		if buffers == nil {
			return
		}
		if err := func() error {
			for _, b := range buffers {
				if err := this.codec.Write(w, b); err != nil {
					return err
				}
			}
			return w.Flush()
		}(); err != nil {
			this.sendQueue.Close()
			return
		}
	}
}

func (this *TCPConnection) Send(b []byte) {
	this.sendQueue.Append(b)
}

func (this *TCPConnection) close() {
	this.conn.SetLinger(0)
	this.conn.Close()
	//this.closeWrite()
	this.sendQueue.Close()
}

func (this *TCPConnection) Shutdown() {
	// this.closeWrite()
	this.sendQueue.Close()
}

func (this *TCPConnection) onConnect() {
	if this.connectFunc != nil {
		this.connectFunc(this)
	}
}

func (this *TCPConnection) onDisonnect() {
	if this.disconnectFunc != nil {
		this.disconnectFunc(this)
	}
}

func (this *TCPConnection) onReceive(b []byte) {
	if this.receiveFunc != nil {
		this.receiveFunc(this, b)
	}
}
