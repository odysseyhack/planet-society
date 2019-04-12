package protocol

import "github.com/odysseyhack/planet-society/protocol/cryptography"

// Message is a message exchanged inside network
type Message struct {
	Header
	Body
}

// Header contains only information needed to route message
type Header struct {
	Source      cryptography.Key32
	Destination cryptography.Key32
	Topic       cryptography.Key32
}

// Body contains payload which is handled by service
type Body struct {
	Payload []byte
}
