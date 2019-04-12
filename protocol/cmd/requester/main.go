package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
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
	fmt.Println("-> received positive pre transaction response")
	fmt.Println("-> sending transaction request")
	fmt.Println("-> received positive transaction response")
}

func connectToResponder() (*transport.Conn, error) {
	u := url.URL{Scheme: "ws", Host: responderAddress, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return &transport.Conn{Conn: conn}, err
}
