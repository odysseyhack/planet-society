package protocol

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/odysseyhack/planet-society/protocol/models"
)

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

type BlockChainExample struct {
}

func (b *BlockChainExample) Validate(request *models.PreTransactionRequest) bool {
	data, err := json.Marshal(request)
	if err != nil {
		return false
	}

	resp, err := http.Client{}.Post("http://block-chain-verification.com:3033", "json", bytes.NewReader(data))
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}
