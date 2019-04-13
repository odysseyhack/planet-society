package utils

import (
	"testing"

	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/models"
)

func TestGetPermissionNo(t *testing.T) {
	if GetPermission("non_exist") == nil {
		t.Errorf("GetPermission returned nil")
	}
}

func TestGetPermission(t *testing.T) {
	transactionID := cryptography.RandomKey32()
	permission := &models.Permission{
		TransactionID: transactionID.String(),
		Description:   "test",
	}

	AddPermission(permission)
	gotPerm := GetPermission(permission.TransactionID)
	if gotPerm == nil {
		t.Errorf("GetPermission returned nil")
	}

	if gotPerm.Description != "test" {
		t.Errorf("GetPermission incorrect permission")
	}
}
