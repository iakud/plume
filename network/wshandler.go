package network

type WSHandler interface {
	Connect(conn *WSConn, connected bool)
	Receive(conn *WSConn, messageType int, data []byte)
}

type defaultWSHandler struct {
}

func (*defaultWSHandler) Connect(*WSConn, bool) {

}

func (*defaultWSHandler) Receive(*WSConn, int, []byte) {

}

var DefaultWSHandler = &defaultWSHandler{}
