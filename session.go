package falcon

import (
	"github.com/iakud/falcon/codec"
	"github.com/iakud/falcon/net"
)

type eventConnect struct {
	session     *Session
	connectFunc func(session *Session)
}

func (this *eventConnect) Run() {
	this.connectFunc(this.session)
}

type eventDisconnect struct {
	session        *Session
	disconnectFunc func(session *Session)
}

func (this *eventDisconnect) Run() {
	this.disconnectFunc(this.session)
}

type eventReceive struct {
	session     *Session
	b           []byte
	receiveFunc func(session *Session, b []byte)
}

func (this *eventReceive) Run() {
	this.receiveFunc(this.session, this.b)
}

type Session struct {
	loop *EventLoop

	codec      codec.Codec
	connection *net.TCPConnection
}

func NewSession(loop *EventLoop, codec codec.Codec, connection *net.TCPConnection) *Session {
	session := &Session{
		loop:       loop,
		codec:      codec,
		connection: connection,
	}
	return session
}

func (this *Session) Serve(connectFunc, disconnectFunc func(session *Session), receiveFunc func(session *Session, b []byte)) {
	this.connection.ServeCodec(this.codec, func(connection *net.TCPConnection) {
		this.loop.RunInLoop(&eventConnect{this, connectFunc})

	}, func(connection *net.TCPConnection) {
		this.loop.RunInLoop(&eventDisconnect{this, disconnectFunc})

	}, func(connection *net.TCPConnection, b []byte) {
		this.loop.RunInLoop(&eventReceive{this, b, receiveFunc})
	})
}

func (this *Session) Send(b []byte) {
	if b == nil {
		return
	}
	this.connection.Send(b)
}

func (this *Session) Shutdown() {
	this.connection.Shutdown()
}
