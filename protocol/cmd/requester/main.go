package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/url"
	"os"

	"github.com/odysseyhack/planet-society/protocol/utils"

	"github.com/gorilla/websocket"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/models"
	"github.com/odysseyhack/planet-society/protocol/protocol"
	"github.com/odysseyhack/planet-society/protocol/transport"
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
	requester        = "John Smith"
)

var (
	keychain           *cryptography.Keychain
	transactionID      *cryptography.Key32
	responderPublicKey *cryptography.Key32
)

func main() {
	fmt.Println("-> generating transaction context")
	if err := createContext(); err != nil {
		fail("-> generating transaction context failed:", err)
	}
	fmt.Println("-> connecting to the responder")

	conn, err := connectToResponder()
	if err != nil {
		fail("-> failed to connect to responder")
	}

	defer func() {
		fmt.Println("-> closing connection to the responder")
		if err := conn.Close(); err != nil {
			fmt.Println("-> failed to close connection to the responder")
		}
	}()

	if err := preTransact(conn); err != nil {
		fmt.Println("-> pre transaction failed:", err)
		os.Exit(1)
	}

	if err := transact(conn); err != nil {
		fmt.Println("-> transaction failed:", err)
		os.Exit(1)
	}
}

func connectToResponder() (*Transport, error) {
	u := url.URL{Scheme: "ws", Host: responderAddress, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("-> connected to the responder")
	return &Transport{&transport.Conn{Conn: conn}}, err
}

func createContext() error {
	k, err := cryptography.OneShotKeychain()
	if err != nil {
		return err
	}

	responderKey, err := utils.ReadKeyFromDir()
	if err != nil {
		return err
	}

	tID := cryptography.RandomKey32()
	keychain = k
	transactionID = &tID
	responderPublicKey = &responderKey

	fmt.Println("-> using responder public key:", responderKey.String())
	fmt.Println("-> using unique transaction ID:", transactionID.String())

	return nil
}

func fail(a ...interface{}) {
	fmt.Println(a)
	os.Exit(1)
}

type Transport struct {
	conn *transport.Conn
}

func (t *Transport) SendMessage(topic cryptography.Key32, payload interface{}) error {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(payload); err != nil {
		return err
	}

	msg := &protocol.Message{
		Header: protocol.Header{
			Source:      keychain.MainPublicKey,
			Destination: *responderPublicKey,
			Topic:       topic,
		},
		Body: protocol.Body{
			Payload: buffer.Bytes(),
		},
	}
	return t.conn.Write(msg)
}

func (t *Transport) Close() error {
	return t.conn.Close()
}

func (t *Transport) Read() (*protocol.Message, error) {
	return t.conn.Read()
}

func preTransact(conn *Transport) error {
	fmt.Println("-> sending pre transaction request")

	if err := conn.SendMessage(
		protocol.TopicPreTransactionRequest,
		&models.PreTransactionRequest{
			TransactionID:      models.Key32{Key: *transactionID},
			SignaturePublicKey: models.Key32{Key: keychain.SignaturePublicKey},
			MainPublicKey:      models.Key32{Key: keychain.MainPublicKey},
			Requester:          "John Smith",
		},
	); err != nil {
		return err
	}

	msg, err := conn.Read()
	if err != nil {
		return err
	}

	var preTransactionReply models.PreTransactionReply
	if err := gob.NewDecoder(bytes.NewBuffer(msg.Body.Payload)).Decode(&preTransactionReply); err != nil {
		return fmt.Errorf("pre transaction reply invalid payload: %s", err)
	}

	if !preTransactionReply.Success {
		return fmt.Errorf("pre transaction was not successful")
	}
	fmt.Println("-> received positive pre transaction response")
	return nil
}

func transact(conn *Transport) error {
	fmt.Println("-> sending transaction request")
	err := conn.SendMessage(
		protocol.TopicTransactionRequest,
		&models.TransactionRequest{
			TransactionID: models.Key32{Key: *transactionID},
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
	msg, err := conn.Read()
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
	return nil
}
