package models

import (
	"encoding/json"
	"testing"
)

func TestNewKeyAnnotation(t *testing.T) {
	ka := NewKeyAnnotation("testuser")
	_, err := json.Marshal(ka)
	if err != nil {
		t.Errorf(err.Error())
	}
}
