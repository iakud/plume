package network

type TCPHandler interface {
	Connected(*TCPConnection)
	Disconnected(*TCPConnection)
	Receive(*TCPConnection, []byte)
}
