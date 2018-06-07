package echo

import (
	"testing"
)

func TestEcho(t *testing.T) {
	echoServer := NewEchoServer("localhost:8000")
	echoServer.Start()

	echoClient := NewEchoClient("localhost:8000")
	echoClient.Start()

	echoClient.Done()
	echoServer.Close()
}
