package protocol

import "github.com/odysseyhack/planet-society/protocol/models"

type AuthorizationPlugin interface {
	Authorize(input *models.PermissionNotificationRequest) (*models.PermissionNotificationResponse, error)
}
