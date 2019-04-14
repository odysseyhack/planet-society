package transport

import (
	"testing"
	"time"

	"github.com/odysseyhack/planet-society/protocol/protocol"
)

func TestNewWebsocket(t *testing.T) {
	connection := make(chan protocol.Conn)

	if NewWebsocket(connection) == nil {
		t.Errorf("NewWebsocket returned nil")
	}
}

func TestWebsocketClose(t *testing.T) {
	connection := make(chan protocol.Conn)
	done := make(chan bool)
	ws := NewWebsocket(connection)

	go func() {
		ws.Listen(":12222")
		done <- true
	}()
	time.Sleep(time.Millisecond * 100)
	if err := ws.Stop(); err != nil {
		t.Errorf("Stop failed: %s", err)
	}
}
