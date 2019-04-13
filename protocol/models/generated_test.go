package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonNotificationRequest(t *testing.T) {
	data, err := json.Marshal(PermissionNotificationResponse{Accepted: true, TransactionID: "10923213981203981"})
	if err != nil {
		t.Errorf("failed to marshal: %s", err)
	}
	fmt.Println(string(data))
}
