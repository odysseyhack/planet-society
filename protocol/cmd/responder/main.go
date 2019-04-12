package main

import (
	"io/ioutil"
	"os"

	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/utils"
)
import log "github.com/sirupsen/logrus"

func main() {
	utils.ConfigureLogger()
	log.Infoln("responder application")

	log.Infoln("creating temporary directory")
	dir, err := ioutil.TempDir("", "responder")
	if err != nil {
		log.Fatalln("failed to create temporary dir:", err)
	}
	log.Infoln("using temporary directory:", dir)

	defer func() {
		log.Infoln("removing temporary directory:", dir)
		if err := os.RemoveAll(dir); err != nil {
			log.Warningln("failed to clean temporary directory:", err)
		}
	}()

	log.Infoln("generating one shot keychain")
	keychain, err := cryptography.OneShotKeychain()
	if err != nil {
		log.Fatalln("failed to generate keychain:", err)
	}
	log.Infoln("keychain main public key:", keychain.MainPublicKey.String())
	log.Infoln("keychain signature key:  ", keychain.SignaturePublicKey.String())

	if err := utils.WriteKeyToDir(keychain.MainPublicKey); err != nil {
		log.Fatalln("failed to write public key:", err)
	}

	defer func() {
		log.Infoln("removing public key")
		if err := utils.CleanKey(); err != nil {
			log.Warningln("cleaning key failed:", err)
		}
	}()

	select {}
}
