package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
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
    birth_date
    email
    BSN
  }
  passport {
    number
    expiration
    country
  }
  
	bankingDetails {
    IBAN
    bank
    nameOnCard
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

func sign() (string, error) {
	signer := cryptography.NewSigner(keychain.SignaturePrivateKey, keychain.SignaturePublicKey)
	signature, err := signer.Sign([]byte(query))
	return hex.EncodeToString(signature), err
}

func createTransactMessage() (*models.TransactionRequest, error) {
	signature, err := sign()
	if err != nil {
		return nil, err
	}

	return &models.TransactionRequest{
		TransactionID: models.Key32{Key: *transactionID},
		Query:         query,
		Title:         "Provide permission for completing",
		Description:   "T-mobile monthly plan(unlimited data), 65 euro, iPhone XR 256GB",
		LawApplying:   "European Union",
		Signature:     signature,
	}, nil
}

func transact(conn *Transport) error {
	smsg, err := createTransactMessage()
	if err != nil {
		return err
	}
	fmt.Println("-> sending transaction request")
	if err := conn.SendMessage(protocol.TopicTransactionRequest, smsg); err != nil {
		return err
	}

	fmt.Println("-> Waiting for transaction reply")
	msg, err := conn.Read()
	if err != nil {
		return err
	}

	var transactionReply models.TransactionReply
	if err := gob.NewDecoder(bytes.NewBuffer(msg.Body.Payload)).Decode(&transactionReply); err != nil {
		return err
	}

	if transactionReply.Error != nil {
		return fmt.Errorf("errors in response: %v", *transactionReply.Error)
	}

	if transactionReply.Content == nil {
		return fmt.Errorf("-> transaction failed: content is nil")
	}
	PrintReply(&transactionReply)
	return nil
}

func PrintReply(reply *models.TransactionReply) {
	fmt.Println("-> received positive transaction response")
	fmt.Println("Transaction content:", *reply.Content)
}
