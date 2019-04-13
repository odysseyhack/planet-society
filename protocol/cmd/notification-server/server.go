package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/odysseyhack/planet-society/protocol/models"
)

type Server struct {
	notification *models.PermissionNotificationRequest
	reply        *models.PermissionNotificationResponse
}

func (s *Server) Listen(addr string) error {
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.createRouter(),
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) createRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/notification-get", s.notificationGet)
	router.HandleFunc("/notification-put", s.notificationPut)
	router.HandleFunc("/reply-put", s.replyPut)
	router.HandleFunc("/reply-get", s.replyGet)
	return router
}

func (s *Server) notificationGet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if s.notification == nil {
		fmt.Println(time.Now().Format(time.RFC3339), "notificationGet: no data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(s.notification)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339), "notificationGet: unmarshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(time.Now().Format(time.RFC3339), "writing notification")
	s.notification = nil
	_, _ = w.Write(data)
}

func (s *Server) notificationPut(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339), "notificationGet: read err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var notification models.PermissionNotificationRequest
	if err := json.Unmarshal(data, &notification); err != nil {
		fmt.Println(time.Now().Format(time.RFC3339), "notificationGet: Unmarshal err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(time.Now().Format(time.RFC3339), "put new notification")
	s.notification = &notification
}

func (s *Server) replyPut(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reply models.PermissionNotificationResponse
	if err := json.Unmarshal(data, &reply); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(time.Now().Format(time.RFC3339), "put new reply")

	s.reply = &reply
}

func (s *Server) replyGet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if s.reply == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(s.reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.reply = nil
	_, _ = w.Write(data)
}
