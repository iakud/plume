package network

type WSHandler interface {
	Connect(conn *WSConn, connected bool)
	Receive(conn *WSConn, messageType WSMessageType, data []byte)
}

type defaultWSHandler struct {
}

func (*defaultWSHandler) Connect(*WSConn, bool) {

}

func (*defaultWSHandler) Receive(*WSConn, WSMessageType, []byte) {

}

var DefaultWSHandler = &defaultWSHandler{}
