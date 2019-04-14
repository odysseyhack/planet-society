package protocol

import (
	"testing"

	"github.com/odysseyhack/planet-society/protocol/cryptography"
)

func TestNewQueue(t *testing.T) {
	if NewQueue() == nil {
		t.Errorf("NewQueue returned nil")
	}
}

func TestNewQueueAdd(t *testing.T) {
	queue := NewQueue()
	entry := &Entry{
		TransactionID: cryptography.RandomKey32(),
	}

	if err := queue.Add(entry); err != nil {
		t.Errorf("queue.Add failed %s", err)
	}
}
