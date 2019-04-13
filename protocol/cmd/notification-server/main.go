package main

import (
	"github.com/odysseyhack/planet-society/protocol/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	utils.ConfigureLogger()
	server := &Server{}
	if err := server.Listen(":80"); err != nil {
		log.Fatalln(err)
	}
}
