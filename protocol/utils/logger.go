package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func ConfigureLogger() {
	// Log as  default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)
}