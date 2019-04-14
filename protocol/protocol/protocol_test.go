package protocol

import "testing"

func TestNewProtocol(t *testing.T) {
	if NewProtocol(nil) == nil {
		t.Errorf("NewProtocol returned nil")
	}
}

func TestProtocolStopLoop(t *testing.T) {
	proto := NewProtocol(nil)
	done := make(chan bool)
	go func() {
		proto.Loop()
		done <- true
	}()
	proto.Stop()
	<-done
}
