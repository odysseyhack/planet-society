package main

import "github.com/odysseyhack/planet-society/protocol/models"

type AlwaysAcceptPlugin struct {
}

func (a *AlwaysAcceptPlugin) Authorize(input *models.PermissionNotificationRequest) (*models.PermissionNotificationResponse, error) {
	return &models.PermissionNotificationResponse{TransactionID: input.TransactionID, Accepted: true}, nil
}
