package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/models"
	"github.com/odysseyhack/planet-society/protocol/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	utils.ConfigureLogger()

	router := mux.NewRouter()
	router.HandleFunc("/notification-get", getHandler)

	server := &http.Server{
		Addr: ":80",
		// todo: review timeouts
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:     router,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln("server failed:",err)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("handling new request")

	transactionID := cryptography.RandomKey32()
	requesterPublicKey := cryptography.RandomKey32()

	personalDetails := models.ItemField{
		Item: "personalDetails",
		Fields: []string{"name", "surname"},
	}

	address := models.ItemField{
		Item: "address",
		Fields: []string{"city", "street","country"},
	}

	notification := models.PermissionNotificationRequest{
		TransactionID: transactionID.String(),
		RequesterPublicKey:requesterPublicKey.String(),
		Reason: "Please provide me your details to finalize transaction",
		Date: time.Now().Format(time.RFC3339),
		Verification: []string{"digid.nl","planet-blockchain"},
		RequesterName: "European Shop inc",
		Item: []models.ItemField{personalDetails,address},
	}

	data, err := json.Marshal(notification)
	if err != nil {
		log.Warningln("failed to serialize:", err)
		return
	}

	w.Write(data)
}


