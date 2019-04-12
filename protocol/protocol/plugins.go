package protocol

import "github.com/odysseyhack/planet-society/protocol/models"

type AuthorizationPlugin interface {
	Authorize(input *models.PermissionNotificationRequest) (*models.PermissionNotificationResponse, error)
}

type Plugins struct {
	preTransactionValidator []PreTransactionValidator
}

type PreTransactionValidator interface {
	Validate(request *models.PreTransactionRequest) bool
}

func (p *Plugins) ValidatePreTransaction(request *models.PreTransactionRequest) bool {
	for i := range p.preTransactionValidator {
		if !p.preTransactionValidator[i].Validate(request) {
			return false
		}
	}
	return true
}
