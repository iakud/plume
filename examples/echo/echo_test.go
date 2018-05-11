package echo

import (
	"log"
	"testing"
)

func TestEcho(t *testing.T) {
	echoServer := NewEchoServer("localhost:8000")
	if err := echoServer.Start(); err != nil {
		log.Fatalln(err)
	}
	echoClient := NewEchoClient("localhost:8000")
	if err := echoClient.Start(); err != nil {
		log.Fatalln(err)
	}
	echoClient.Done()
	echoServer.Close()
}
