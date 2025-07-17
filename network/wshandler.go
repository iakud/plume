package network

type WSHandler interface {
	Connect(conn *WSConnection, connected bool)
	Receive(conn *WSConnection, messageType int, data []byte)
}

type defaultWSHandler struct {
}

func (*defaultWSHandler) Connect(*WSConnection, bool) {

}

func (*defaultWSHandler) Receive(*WSConnection, int, []byte) {

}

var DefaultWSHandler = &defaultWSHandler{}
