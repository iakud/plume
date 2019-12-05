package network

type TCPHandler interface {
	Connect(conn *TCPConnection, connected bool)
	Receive(conn *TCPConnection, buf []byte)
}

type defaultTCPHandler struct {
}

func (*defaultTCPHandler) Connect(*TCPConnection, bool) {

}

func (*defaultTCPHandler) Receive(*TCPConnection, []byte) {

}

var DefaultTCPHandler *defaultTCPHandler = &defaultTCPHandler{}
