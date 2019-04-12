package protocol

import "github.com/odysseyhack/planet-society/protocol/cryptography"

var (
	TopicPreTransactionRequest = cryptography.Key32{'1'}
	TopicPreTransactionReply   = cryptography.Key32{'2'}
	TopicTransactionRequest    = cryptography.Key32{'3'}
	TopicTransactionReply      = cryptography.Key32{'3'}
)
