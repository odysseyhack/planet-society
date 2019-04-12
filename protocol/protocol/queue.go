package protocol

import (
	"sync"

	"github.com/odysseyhack/planet-society/protocol/cryptography"
)

type Entry struct {
	TransactionID      cryptography.Key32
	RequesterName      string
	Authorization      map[string]string
	RequesterPublicKey cryptography.Key32
}

type Queue struct {
	sync.RWMutex
	entries map[cryptography.Key32]*Entry
}

func NewQueue() *Queue {
	return &Queue{
		entries: make(map[cryptography.Key32]*Entry),
	}
}

func (q *Queue) Add(e *Entry) error {
	q.Lock()
	defer q.Unlock()

	if _, ok := q.entries[e.TransactionID]; ok {
		return ErrIDAlreadyUsed
	}
	q.entries[e.TransactionID] = e
	return nil
}

func (q *Queue) Get(TransactionID cryptography.Key32) (*Entry, bool) {
	entry, ok := q.entries[TransactionID]
	return entry, ok
}
