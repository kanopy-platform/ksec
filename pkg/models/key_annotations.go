package models

// KeyAnnotation holds metadata about individual Secrets keys
type KeyAnnotation struct {
	UpdatedBy   string `json:"updatedBy"`
	LastUpdated string `json:"lastUpdated"`
}
