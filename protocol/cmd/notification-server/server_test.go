package main

import (
	"testing"
	"time"
)

func TestServerClose(t *testing.T) {
	server := Server{}
	done := make(chan bool)
	go func() {
		server.Listen(":12121")
		done <- true
	}()
	time.Sleep(time.Millisecond * 100)
	if err := server.Stop(); err != nil {
		t.Errorf("Stop failed: %s", err)
	}
	<-done
}

func TestServerNewRouter(t *testing.T) {
	server := Server{}
	if server.createRouter() == nil {
		t.Errorf("router creation failed")
	}
}
