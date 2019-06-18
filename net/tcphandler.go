package net

type TCPHandler interface {
	Connect(*TCPConnection)
	Disconnect(*TCPConnection)
	Receive(*TCPConnection, []byte)
}

type defaultTCPHandler struct {
}

func (*defaultTCPHandler) Connect(*TCPConnection) {

}

func (*defaultTCPHandler) Disconnect(*TCPConnection) {

}

func (*defaultTCPHandler) Receive(*TCPConnection, []byte) {

}

var DefaultTCPHandler *defaultTCPHandler = &defaultTCPHandler{}
