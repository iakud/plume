package echo

import (
	"log"
	"testing"
)

func TestEcho(t *testing.T) {
	echoServer := NewEchoServer("127.0.0.1:9000")
	if err := echoServer.Start(); err != nil {
		log.Fatalln(err)
	}
	echoClient := NewEchoClient("127.0.0.1:9000")
	if err := echoClient.Start(); err != nil {
		log.Fatalln(err)
	}
	<-echoClient.Done()
	echoServer.Close()
}
