package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/odysseyhack/planet-society/protocol/models"
)

type AlwaysAcceptPlugin struct {
}

func (a *AlwaysAcceptPlugin) Authorize(input *models.PermissionNotificationRequest) (*models.PermissionNotificationResponse, error) {
	return &models.PermissionNotificationResponse{TransactionID: input.TransactionID, Accepted: true}, nil
}

const (
	notificationPutAddr = "http://51.15.52.136/notification-put"
	replyGetAddr        = "http://51.15.52.136/reply-get"
)

type IOSPlugin struct {
}

func (i *IOSPlugin) Authorize(input *models.PermissionNotificationRequest) (*models.PermissionNotificationResponse, error) {
	if err := sendPermissionNotificationRequest(input); err != nil {
		return nil, err
	}
	return responseGetLoop()
}

func responseGetLoop() (*models.PermissionNotificationResponse, error) {
	for i := 0; i < 120; i++ {
		resp, err := getPermissionNotificationResponse()
		if err == nil {
			return resp, nil
		}
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("timeout")
}

func getPermissionNotificationResponse() (*models.PermissionNotificationResponse, error) {
	client := http.Client{}
	resp, err := client.Get(replyGetAddr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sendPermissionNotificationRequest: got %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var notificationResponse models.PermissionNotificationResponse
	if err := json.Unmarshal(data, &notificationResponse); err != nil {
		return nil, err
	}
	return &notificationResponse, err
}

func sendPermissionNotificationRequest(input *models.PermissionNotificationRequest) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	client := http.Client{}

	rq, err := http.NewRequest(http.MethodGet, notificationPutAddr, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	resp, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sendPermissionNotificationRequest: got %v", resp.StatusCode)
	}

	return nil
}
