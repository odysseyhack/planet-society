package utils

import "github.com/odysseyhack/planet-society/protocol/models"

var (
	permMap = make(map[string]*models.Permission)
)

func AddPermission(perm *models.Permission) {
	permMap[perm.TransactionID] = perm
}

func GetPermission(transactionID string) *models.Permission {
	perm, ok := permMap[transactionID]
	if ok {
		return perm
	}
	return generatedPermission(transactionID)
}

func generatedPermission(transactionID string) *models.Permission {
	return &models.Permission{TransactionID: transactionID}
}
