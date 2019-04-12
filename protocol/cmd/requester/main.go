package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/models"
	"github.com/odysseyhack/planet-society/protocol/protocol"
	"github.com/odysseyhack/planet-society/protocol/transport"
	"github.com/odysseyhack/planet-society/protocol/utils"
)

const query = `
query {
  personalDetails {
    name
    surname
  }
  address {
    country
    city
    street
  }
}
`

const (
	responderAddress = ":15000"
)

func main() {
	fmt.Println("-> generating identity")
	keychain, err := cryptography.OneShotKeychain()
	if err != nil {
		fmt.Println("-> failed to generate identity")
		os.Exit(1)
	}

	transactionID := cryptography.RandomKey32()
	fmt.Println("-> using unique transaction ID:", transactionID.String())
	fmt.Println("-> connecting to the responder")

	responderKey, err := utils.ReadKeyFromDir()
	if err != nil {
		fmt.Println("-> failed to get responder public key")
		os.Exit(1)
	}
	fmt.Println("-> using responder public key:", responderKey.String())

	conn, err := connectToResponder()
	if err != nil {
		fmt.Println("-> failed to connect to responder")
		os.Exit(1)
	}

	defer func() {
		fmt.Println("-> closing connection to the responder")
		if err := conn.Close(); err != nil {
			fmt.Println("-> failed to close connection to the responder")
		}
	}()
	fmt.Println("-> connected to the responder")
	fmt.Println("-> sending pre transaction request")

	sendMsg := func(topic cryptography.Key32, payload interface{}) error {
		var buffer bytes.Buffer

		if err := gob.NewEncoder(&buffer).Encode(payload); err != nil {
			return err
		}

		msg := &protocol.Message{
			Header: protocol.Header{
				Source:      keychain.MainPublicKey,
				Destination: responderKey,
				Topic:       topic,
			},
			Body: protocol.Body{
				Payload: buffer.Bytes(),
			},
		}
		return conn.Write(msg)
	}

	err = sendMsg(
		protocol.TopicPreTransactionRequest,
		&models.PreTransactionRequest{
			TransactionID:      models.Key32{Key: transactionID},
			SignaturePublicKey: models.Key32{Key: keychain.SignaturePublicKey},
			MainPublicKey:      models.Key32{Key: keychain.MainPublicKey},
			Requester:          "John Smith",
		},
	)
	if err != nil {
		fmt.Println("-> sending pre transaction failed:", err)
		os.Exit(1)
	}

	msg, err := conn.Read()
	if err != nil {
		fmt.Println("-> read failed:", err)
		os.Exit(1)
	}

	var preTransactionReply models.PreTransactionReply
	if err := gob.NewDecoder(bytes.NewBuffer(msg.Body.Payload)).Decode(&preTransactionReply); err != nil {
		fmt.Println("-> pre transaction reply: invalid payload:", err)
		os.Exit(1)
	}

	if !preTransactionReply.Success {
		fmt.Println("-> pre transaction was not successful")
		os.Exit(1)
	}

	fmt.Println("-> received positive pre transaction response")
	fmt.Println("-> sending transaction request")
	err = sendMsg(
		protocol.TopicTransactionRequest,
		&models.TransactionRequest{
			TransactionID: models.Key32{Key: transactionID},
			Query:         query,
			Reason:        "please provide me following details to finalize shipment",
			LawApplying:   "European Union",
		},
	)
	if err != nil {
		fmt.Println("-> sending transaction request failed", err)
		os.Exit(1)
	}

	fmt.Println("-> Waiting for transaction reply")
	msg, err = conn.Read()
	if err != nil {
		fmt.Println("-> Read failed", err)
		os.Exit(1)
	}

	var transactionReply models.TransactionReply
	if err := gob.NewDecoder(bytes.NewBuffer(msg.Body.Payload)).Decode(&transactionReply); err != nil {
		// invalid data, cannot reply
		fmt.Println("-> transaction reply: invalid payload:", err)
		os.Exit(1)
	}

	if transactionReply.Error != nil {
		fmt.Println("-> transaction failed:", *transactionReply.Error)
		os.Exit(1)
	}

	if transactionReply.Content == nil {
		fmt.Println("-> transaction failed: content is nil")
		os.Exit(1)
	}

	fmt.Println("-> received positive transaction response")
	fmt.Println("Transaction content:", *transactionReply.Content)
}

func connectToResponder() (*transport.Conn, error) {
	u := url.URL{Scheme: "ws", Host: responderAddress, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return &transport.Conn{Conn: conn}, err
}
