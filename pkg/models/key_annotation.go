package models

import (
	"time"
)

const annotationPrefix = "ksec.io"

// KeyAnnotation holds metadata about individual Secrets keys
type KeyAnnotation struct {
	UpdatedBy   string `json:"updatedBy"`
	LastUpdated string `json:"lastUpdated"`
}

// NewKeyAnnotation constructor
func NewKeyAnnotation(authInfo string) *KeyAnnotation {
	return &KeyAnnotation{
		UpdatedBy:   authInfo,
		LastUpdated: time.Now().Format(time.RFC3339),
	}
}
