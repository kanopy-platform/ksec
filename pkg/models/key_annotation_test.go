package models

import (
	"fmt"
	"encoding/json"
	"testing"
)

func TestNewKeyAnnotation(t *testing.T) {
	ka := NewKeyAnnotation("testuser")
	jsonBytes, err := json.Marshal(ka)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(jsonBytes))
}
