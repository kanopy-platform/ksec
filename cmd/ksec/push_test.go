package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanSecretData(t *testing.T) {
	reader := strings.NewReader(`key=value`)
	data := map[string][]byte{}

	assert.NoError(t, scanSecretData(reader, data))
	assert.Equal(t, "value", string(data["key"]))
}
