package main

import (
	"fmt"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
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
	utils.ConfigureLogger()


	transactionID := cryptography.RandomKey32()
	fmt.Println("-> using unique transaction ID:", transactionID.String())
	fmt.Println("-> sending pre transaction request")
	fmt.Println("-> received positive pre transaction response")
	fmt.Println("-> sending transaction request")
	fmt.Println("-> received positive transaction response")
}
