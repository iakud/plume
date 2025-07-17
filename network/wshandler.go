package network

type WSHandler interface {
	Connect(conn *WSConnection, connected bool)
	Receive(conn *WSConnection, data []byte)
}

type defaultWSHandler struct {
}

func (*defaultWSHandler) Connect(*WSConnection, bool) {

}

func (*defaultWSHandler) Receive(*WSConnection, []byte) {

}

var DefaultWSHandler = &defaultWSHandler{}
